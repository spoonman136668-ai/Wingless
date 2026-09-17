#include "ggml.h"
#include "ggml-backend.h"
#include "llama.h"

#include <algorithm>
#include <atomic>
#include <cctype>
#include <clocale>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <fstream>
#include <limits>
#include <mutex>
#include <string>
#include <vector>

static constexpr const char * k_trace_schema = "wingless.nr1.router-trace.v1";
static constexpr const char * k_llama_commit = "5266f24da75dc449bd56cbed7addb9c8e4a6a73e";
static constexpr const char * k_topk_marker = "ffn_moe_topk-";

struct options {
    std::string model_path;
    std::string trace_path;
    std::string prompt_file;
    std::string user_file;
    std::string system_file;
    std::string prompt;
    std::string user;
    std::string system;
    std::string session_id;
    std::string workload_id;
    std::string task_family;
    bool chat_mode = false;
    int n_gpu_layers = 32;
    int n_ctx = 4096;
    int n_batch = 256;
    int n_predict = 128;
};

static void usage(const char * argv0) {
    std::fprintf(stderr,
        "usage: %s -m MODEL --trace TRACE --session ID --workload ID --task FAMILY "
        "(--user-file FILE [--system-file FILE] | --prompt-file FILE | PROMPT) "
        "[-ngl 32] [-c 4096] [-b 256] [-n 128]\n",
        argv0);
}

static bool safe_id(const std::string & value) {
    if (value.empty() || value.size() > 128 || !std::isalnum(static_cast<unsigned char>(value[0]))) {
        return false;
    }
    for (unsigned char c : value) {
        if (!(std::isalnum(c) || c == '_' || c == '.' || c == ':' || c == '-')) {
            return false;
        }
    }
    return true;
}

static bool read_text(const std::string & path, std::string & out) {
    std::ifstream in(path, std::ios::binary);
    if (!in) {
        return false;
    }
    in.seekg(0, std::ios::end);
    const std::streamoff size = in.tellg();
    if (size < 0 || size > 16 * 1024 * 1024) {
        return false;
    }
    in.seekg(0, std::ios::beg);
    out.resize(static_cast<size_t>(size));
    if (size > 0) {
        in.read(out.data(), size);
    }
    return in.good() || in.eof();
}

static bool parse_int(const char * text, int & out) {
    if (text == nullptr || *text == '\0') {
        return false;
    }
    char * end = nullptr;
    const long value = std::strtol(text, &end, 10);
    if (*end != '\0' || value < 0 || value > 1'000'000) {
        return false;
    }
    out = static_cast<int>(value);
    return true;
}

static bool parse_options(int argc, char ** argv, options & opt) {
    for (int i = 1; i < argc; ++i) {
        const std::string arg = argv[i];
        auto next = [&]() -> const char * {
            if (i + 1 >= argc) {
                return nullptr;
            }
            return argv[++i];
        };

        if (arg == "-m" || arg == "--model") {
            const char * v = next(); if (!v) return false; opt.model_path = v;
        } else if (arg == "--trace") {
            const char * v = next(); if (!v) return false; opt.trace_path = v;
        } else if (arg == "--session") {
            const char * v = next(); if (!v) return false; opt.session_id = v;
        } else if (arg == "--workload") {
            const char * v = next(); if (!v) return false; opt.workload_id = v;
        } else if (arg == "--task") {
            const char * v = next(); if (!v) return false; opt.task_family = v;
        } else if (arg == "--prompt-file") {
            const char * v = next(); if (!v) return false; opt.prompt_file = v;
        } else if (arg == "--user-file") {
            const char * v = next(); if (!v) return false; opt.user_file = v;
        } else if (arg == "--system-file") {
            const char * v = next(); if (!v) return false; opt.system_file = v;
        } else if (arg == "-ngl") {
            const char * v = next(); if (!v || !parse_int(v, opt.n_gpu_layers)) return false;
        } else if (arg == "-c") {
            const char * v = next(); if (!v || !parse_int(v, opt.n_ctx)) return false;
        } else if (arg == "-b") {
            const char * v = next(); if (!v || !parse_int(v, opt.n_batch)) return false;
        } else if (arg == "-n") {
            const char * v = next(); if (!v || !parse_int(v, opt.n_predict)) return false;
        } else if (!arg.empty() && arg[0] == '-') {
            return false;
        } else {
            if (!opt.prompt.empty()) opt.prompt.push_back(' ');
            opt.prompt += arg;
        }
    }

    if (opt.model_path.empty() || opt.trace_path.empty() ||
        !safe_id(opt.session_id) || !safe_id(opt.workload_id) || !safe_id(opt.task_family) ||
        opt.n_ctx < 2 || opt.n_batch < 1 || opt.n_predict < 0 || opt.n_gpu_layers < 0) {
        return false;
    }

    const bool has_raw_file = !opt.prompt_file.empty();
    const bool has_raw_text = !opt.prompt.empty();
    const bool has_user_file = !opt.user_file.empty();
    const int input_modes = (has_raw_file ? 1 : 0) + (has_raw_text ? 1 : 0) + (has_user_file ? 1 : 0);
    if (input_modes != 1 || (!opt.system_file.empty() && !has_user_file)) {
        return false;
    }

    if (has_user_file) {
        if (!read_text(opt.user_file, opt.user) || opt.user.empty()) {
            return false;
        }
        if (!opt.system_file.empty() && !read_text(opt.system_file, opt.system)) {
            return false;
        }
        opt.chat_mode = true;
    } else if (has_raw_file) {
        if (!read_text(opt.prompt_file, opt.prompt) || opt.prompt.empty()) {
            return false;
        }
    }
    return true;
}

static bool apply_chat_template(const llama_model * model, options & opt) {
    const char * tmpl = llama_model_chat_template(model, nullptr);
    if (tmpl == nullptr || *tmpl == '\0') {
        std::fprintf(stderr, "model has no default chat template\n");
        return false;
    }

    std::vector<llama_chat_message> messages;
    if (!opt.system.empty()) {
        messages.push_back({"system", opt.system.c_str()});
    }
    messages.push_back({"user", opt.user.c_str()});

    int required = llama_chat_apply_template(tmpl, messages.data(), messages.size(), true, nullptr, 0);
    if (required < 0 || required > 16 * 1024 * 1024) {
        std::fprintf(stderr, "failed to size model chat template output\n");
        return false;
    }

    std::vector<char> formatted(static_cast<size_t>(required) + 1, '\0');
    const int written = llama_chat_apply_template(
        tmpl, messages.data(), messages.size(), true, formatted.data(), formatted.size());
    if (written < 0 || written > required) {
        std::fprintf(stderr, "failed to apply model chat template\n");
        return false;
    }
    opt.prompt.assign(formatted.data(), static_cast<size_t>(written));
    return !opt.prompt.empty();
}

static bool parse_layer(const char * raw_name, int & layer) {
    if (raw_name == nullptr) {
        return false;
    }
    const std::string name(raw_name);
    const size_t pos = name.rfind(k_topk_marker);
    if (pos == std::string::npos) {
        return false;
    }
    size_t p = pos + std::strlen(k_topk_marker);
    if (p >= name.size() || !std::isdigit(static_cast<unsigned char>(name[p]))) {
        return false;
    }
    long value = 0;
    while (p < name.size() && std::isdigit(static_cast<unsigned char>(name[p]))) {
        value = value * 10 + (name[p] - '0');
        if (value > 255) {
            return false;
        }
        ++p;
    }
    layer = static_cast<int>(value);
    return true;
}

static bool read_topk_ids(const ggml_tensor * t, std::vector<int32_t> & ids, const char *& error) {
    if (t == nullptr || t->type != GGML_TYPE_I32 ||
        t->ne[0] < 1 || t->ne[0] > 128 || t->ne[1] < 1 || t->ne[2] != 1 || t->ne[3] != 1) {
        error = "unexpected ffn_moe_topk tensor shape/type";
        return false;
    }
    if (t->nb[0] != sizeof(int32_t)) {
        error = "ffn_moe_topk dimension-0 is not dense I32";
        return false;
    }

    const size_t cols = static_cast<size_t>(t->ne[0]);
    const size_t rows = static_cast<size_t>(t->ne[1]);
    if (cols > std::numeric_limits<size_t>::max() / sizeof(int32_t)) {
        error = "ffn_moe_topk row byte size overflow";
        return false;
    }
    const size_t row_bytes = cols * sizeof(int32_t);
    if (rows > std::numeric_limits<size_t>::max() / row_bytes) {
        error = "ffn_moe_topk dense byte size overflow";
        return false;
    }
    const size_t dense_bytes = rows * row_bytes;
    const size_t count = rows * cols;

    if (static_cast<size_t>(ggml_nelements(t)) != count) {
        error = "unexpected ffn_moe_topk tensor rank";
        return false;
    }

    ids.resize(count);
    if (ggml_is_contiguous(t)) {
        ggml_backend_tensor_get(t, ids.data(), 0, dense_bytes);
        return true;
    }

    // argsort_top_k is logically [n_expert_used, n_tokens], but backend graph
    // rewrites may leave the token rows strided. Copy only the dense I32 row
    // payload into our packed host buffer; never read ggml_nbytes(t) into a
    // ggml_nelements(t)-sized vector because a view can include stride gaps.
    if (t->nb[1] < row_bytes) {
        error = "ffn_moe_topk row stride is smaller than dense row";
        return false;
    }
    ggml_backend_tensor_get_2d(
        t,
        ids.data(),
        0,
        row_bytes,
        rows,
        t->nb[1],
        row_bytes);
    return true;
}

class trace_collector {
public:
    trace_collector(const options & opt) :
        out_(opt.trace_path, std::ios::binary | std::ios::trunc),
        session_(opt.session_id), workload_(opt.workload_id), task_(opt.task_family) {
        pending_.reserve(1 << 20);
    }

    bool ready() const { return out_.is_open(); }
    bool failed() const { return failed_.load(); }

    void set_token_base(int64_t base) {
        std::lock_guard<std::mutex> lock(mu_);
        token_base_ = base;
    }

    bool wants(const ggml_tensor * t) const {
        int layer = -1;
        return t != nullptr && parse_layer(t->name, layer);
    }

    bool collect(ggml_tensor * t) {
        int layer = -1;
        if (t == nullptr || !parse_layer(t->name, layer)) {
            return true;
        }

        std::vector<int32_t> ids;
        const char * read_error = nullptr;
        if (!read_topk_ids(t, ids, read_error)) {
            fail(read_error != nullptr ? read_error : "failed to read ffn_moe_topk tensor");
            return false;
        }

        std::lock_guard<std::mutex> lock(mu_);
        for (int64_t token = 0; token < t->ne[1]; ++token) {
            pending_ += "{\"schema\":\"";
            pending_ += k_trace_schema;
            pending_ += "\",\"session_id\":\"";
            pending_ += session_;
            pending_ += "\",\"workload_id\":\"";
            pending_ += workload_;
            pending_ += "\",\"task_family\":\"";
            pending_ += task_;
            pending_ += "\",\"token_index\":";
            pending_ += std::to_string(token_base_ + token);
            pending_ += ",\"layer\":";
            pending_ += std::to_string(layer);
            pending_ += ",\"experts\":[";
            for (int64_t k = 0; k < t->ne[0]; ++k) {
                if (k != 0) pending_.push_back(',');
                const int32_t id = ids[static_cast<size_t>(token * t->ne[0] + k)];
                if (id < 0 || id > 1023) {
                    fail_locked("expert id out of range");
                    return false;
                }
                pending_ += std::to_string(id);
            }
            pending_ += "]}\n";
        }
        return true;
    }

    bool flush() {
        std::lock_guard<std::mutex> lock(mu_);
        if (failed_.load()) {
            return false;
        }
        out_.write(pending_.data(), static_cast<std::streamsize>(pending_.size()));
        pending_.clear();
        out_.flush();
        if (!out_) {
            fail_locked("trace write failed");
            return false;
        }
        return true;
    }

private:
    void fail(const char * message) {
        std::lock_guard<std::mutex> lock(mu_);
        fail_locked(message);
    }

    void fail_locked(const char * message) {
        if (!failed_.exchange(true)) {
            std::fprintf(stderr, "nr1 trace error: %s\n", message);
        }
    }

    std::ofstream out_;
    std::string session_;
    std::string workload_;
    std::string task_;
    std::string pending_;
    int64_t token_base_ = 0;
    std::atomic<bool> failed_{false};
    mutable std::mutex mu_;
};

static bool trace_eval_cb(ggml_tensor * t, bool ask, void * user_data) {
    auto * collector = static_cast<trace_collector *>(user_data);
    if (collector == nullptr) {
        return false;
    }
    if (ask) {
        return collector->wants(t);
    }
    return collector->collect(t);
}

int main(int argc, char ** argv) {
    std::setlocale(LC_NUMERIC, "C");

    options opt;
    if (!parse_options(argc, argv, opt)) {
        usage(argv[0]);
        return 2;
    }

    std::fprintf(stderr, "NR-1A trace producer: llama.cpp commit %s\n", k_llama_commit);

    trace_collector collector(opt);
    if (!collector.ready()) {
        std::fprintf(stderr, "failed to open trace output: %s\n", opt.trace_path.c_str());
        return 3;
    }

    ggml_backend_load_all();

    llama_model_params model_params = llama_model_default_params();
    model_params.n_gpu_layers = opt.n_gpu_layers;

    llama_model * model = llama_model_load_from_file(opt.model_path.c_str(), model_params);
    if (model == nullptr) {
        std::fprintf(stderr, "failed to load model\n");
        return 4;
    }

    if (opt.chat_mode && !apply_chat_template(model, opt)) {
        llama_model_free(model);
        return 5;
    }

    const llama_vocab * vocab = llama_model_get_vocab(model);
    const bool add_special = opt.chat_mode ? true : llama_vocab_get_add_bos(vocab);
    const int n_prompt = -llama_tokenize(vocab, opt.prompt.c_str(), opt.prompt.size(), nullptr, 0, add_special, true);
    if (n_prompt <= 0 || n_prompt >= opt.n_ctx) {
        std::fprintf(stderr, "prompt token count %d does not fit context %d\n", n_prompt, opt.n_ctx);
        llama_model_free(model);
        return 6;
    }

    std::vector<llama_token> prompt_tokens(static_cast<size_t>(n_prompt));
    if (llama_tokenize(vocab, opt.prompt.c_str(), opt.prompt.size(), prompt_tokens.data(), prompt_tokens.size(), add_special, true) < 0) {
        std::fprintf(stderr, "failed to tokenize prompt\n");
        llama_model_free(model);
        return 7;
    }

    if (n_prompt + opt.n_predict > opt.n_ctx) {
        std::fprintf(stderr, "prompt + prediction exceeds context: %d + %d > %d\n", n_prompt, opt.n_predict, opt.n_ctx);
        llama_model_free(model);
        return 8;
    }

    llama_context_params ctx_params = llama_context_default_params();
    ctx_params.n_ctx = static_cast<uint32_t>(opt.n_ctx);
    ctx_params.n_batch = static_cast<uint32_t>(std::min(opt.n_batch, opt.n_ctx));
    ctx_params.n_ubatch = ctx_params.n_batch;
    ctx_params.no_perf = false;
    ctx_params.cb_eval = trace_eval_cb;
    ctx_params.cb_eval_user_data = &collector;

    llama_context * ctx = llama_init_from_model(model, ctx_params);
    if (ctx == nullptr) {
        std::fprintf(stderr, "failed to create context\n");
        llama_model_free(model);
        return 9;
    }

    auto sparams = llama_sampler_chain_default_params();
    sparams.no_perf = false;
    llama_sampler * sampler = llama_sampler_chain_init(sparams);
    llama_sampler_chain_add(sampler, llama_sampler_init_greedy());

    int64_t n_pos = 0;
    while (n_pos < n_prompt) {
        const int n = std::min<int64_t>(opt.n_batch, n_prompt - n_pos);
        llama_batch batch = llama_batch_get_one(prompt_tokens.data() + n_pos, n);
        collector.set_token_base(n_pos);
        if (llama_decode(ctx, batch) != 0 || collector.failed() || !collector.flush()) {
            std::fprintf(stderr, "prompt decode/trace failed at token %lld\n", static_cast<long long>(n_pos));
            llama_sampler_free(sampler);
            llama_free(ctx);
            llama_model_free(model);
            return 10;
        }
        n_pos += n;
    }

    int generated = 0;
    while (generated < opt.n_predict) {
        const llama_token token = llama_sampler_sample(sampler, ctx, -1);
        if (llama_vocab_is_eog(vocab, token)) {
            break;
        }

        llama_batch batch = llama_batch_get_one(&token, 1);
        collector.set_token_base(n_pos);
        if (llama_decode(ctx, batch) != 0 || collector.failed() || !collector.flush()) {
            std::fprintf(stderr, "generation decode/trace failed at token %lld\n", static_cast<long long>(n_pos));
            llama_sampler_free(sampler);
            llama_free(ctx);
            llama_model_free(model);
            return 11;
        }
        ++n_pos;
        ++generated;
    }

    std::fprintf(stderr, "NR-1A trace complete: mode=%s prompt_tokens=%d generated_tokens=%d\n",
        opt.chat_mode ? "model_chat_template" : "raw_prompt", n_prompt, generated);
    llama_perf_context_print(ctx);

    llama_sampler_free(sampler);
    llama_free(ctx);
    llama_model_free(model);
    llama_backend_free();
    return 0;
}
