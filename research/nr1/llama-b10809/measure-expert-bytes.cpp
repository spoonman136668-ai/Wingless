#include "ggml.h"
#include "gguf.h"

#include <cinttypes>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <limits>
#include <string>
#include <vector>

static constexpr int k_layers = 48;
static constexpr int k_experts = 128;
static constexpr const char * k_schema = "wingless.nr1.expert-storage.v1";
static constexpr const char * k_llama_commit = "5266f24da75dc449bd56cbed7addb9c8e4a6a73e";

struct tensor_info {
    std::string name;
    std::string type;
    uint64_t bytes = 0;
    uint64_t bytes_per_expert = 0;
};

struct layer_info {
    int layer = 0;
    tensor_info gate;
    tensor_info up;
    tensor_info down;
    uint64_t bytes = 0;
    uint64_t bytes_per_expert = 0;
};

static bool checked_add(uint64_t a, uint64_t b, uint64_t & out) {
    if (a > std::numeric_limits<uint64_t>::max() - b) {
        return false;
    }
    out = a + b;
    return true;
}

static bool inspect_tensor(
        const gguf_context * ctx,
        const std::string & name,
        tensor_info & out) {
    const int64_t id = gguf_find_tensor(ctx, name.c_str());
    if (id < 0) {
        std::fprintf(stderr, "missing expert tensor: %s\n", name.c_str());
        return false;
    }

    const int64_t * ne = gguf_get_tensor_ne(ctx, id);
    if (ne == nullptr || ne[0] <= 0 || ne[1] <= 0 || ne[2] != k_experts || ne[3] != 1) {
        std::fprintf(stderr, "unexpected expert tensor shape for %s: (%" PRId64 ", %" PRId64 ", %" PRId64 ", %" PRId64 ")\n",
            name.c_str(), ne ? ne[0] : -1, ne ? ne[1] : -1, ne ? ne[2] : -1, ne ? ne[3] : -1);
        return false;
    }

    const size_t size = gguf_get_tensor_size(ctx, id);
    if (size == 0 || size % k_experts != 0) {
        std::fprintf(stderr, "expert tensor size is not evenly divisible by %d: %s size=%zu\n",
            k_experts, name.c_str(), size);
        return false;
    }

    out.name = name;
    out.type = ggml_type_name(gguf_get_tensor_type(ctx, id));
    out.bytes = static_cast<uint64_t>(size);
    out.bytes_per_expert = static_cast<uint64_t>(size / k_experts);
    return true;
}

static void print_json_string(const std::string & value) {
    std::putchar('"');
    for (unsigned char c : value) {
        switch (c) {
            case '"': std::fputs("\\\"", stdout); break;
            case '\\': std::fputs("\\\\", stdout); break;
            case '\b': std::fputs("\\b", stdout); break;
            case '\f': std::fputs("\\f", stdout); break;
            case '\n': std::fputs("\\n", stdout); break;
            case '\r': std::fputs("\\r", stdout); break;
            case '\t': std::fputs("\\t", stdout); break;
            default:
                if (c < 0x20) {
                    std::printf("\\u%04x", static_cast<unsigned int>(c));
                } else {
                    std::putchar(c);
                }
        }
    }
    std::putchar('"');
}

static void print_tensor(const char * key, const tensor_info & t) {
    std::printf("\"%s\":{\"name\":", key);
    print_json_string(t.name);
    std::fputs(",\"type\":", stdout);
    print_json_string(t.type);
    std::printf(",\"bytes\":%" PRIu64 ",\"bytes_per_expert\":%" PRIu64 "}", t.bytes, t.bytes_per_expert);
}

int main(int argc, char ** argv) {
    if (argc != 2 || argv[1] == nullptr || argv[1][0] == '\0') {
        std::fprintf(stderr, "usage: %s MODEL.gguf\n", argv[0]);
        return 2;
    }

    gguf_init_params params = {
        /*.no_alloc =*/ false,
        /*.ctx      =*/ nullptr,
    };
    gguf_context * ctx = gguf_init_from_file(argv[1], params);
    if (ctx == nullptr) {
        std::fprintf(stderr, "failed to read GGUF metadata: %s\n", argv[1]);
        return 3;
    }

    uint64_t total_tensor_bytes = 0;
    const int64_t n_tensors = gguf_get_n_tensors(ctx);
    for (int64_t i = 0; i < n_tensors; ++i) {
        const uint64_t size = static_cast<uint64_t>(gguf_get_tensor_size(ctx, i));
        if (!checked_add(total_tensor_bytes, size, total_tensor_bytes)) {
            std::fprintf(stderr, "total tensor byte count overflow\n");
            gguf_free(ctx);
            return 4;
        }
    }

    std::vector<layer_info> layers;
    layers.reserve(k_layers);
    uint64_t pool_bytes = 0;
    uint64_t min_bytes_per_expert = std::numeric_limits<uint64_t>::max();
    uint64_t max_bytes_per_expert = 0;

    for (int layer = 0; layer < k_layers; ++layer) {
        layer_info info;
        info.layer = layer;
        const std::string prefix = "blk." + std::to_string(layer);
        if (!inspect_tensor(ctx, prefix + ".ffn_gate_exps.weight", info.gate) ||
            !inspect_tensor(ctx, prefix + ".ffn_up_exps.weight", info.up) ||
            !inspect_tensor(ctx, prefix + ".ffn_down_exps.weight", info.down)) {
            gguf_free(ctx);
            return 5;
        }

        uint64_t tmp = 0;
        if (!checked_add(info.gate.bytes, info.up.bytes, tmp) ||
            !checked_add(tmp, info.down.bytes, info.bytes) ||
            !checked_add(pool_bytes, info.bytes, pool_bytes)) {
            std::fprintf(stderr, "expert storage byte count overflow\n");
            gguf_free(ctx);
            return 6;
        }
        info.bytes_per_expert = info.bytes / k_experts;
        if (info.bytes % k_experts != 0) {
            std::fprintf(stderr, "layer %d expert storage is not evenly divisible by %d\n", layer, k_experts);
            gguf_free(ctx);
            return 7;
        }
        if (info.bytes_per_expert < min_bytes_per_expert) min_bytes_per_expert = info.bytes_per_expert;
        if (info.bytes_per_expert > max_bytes_per_expert) max_bytes_per_expert = info.bytes_per_expert;
        layers.push_back(info);
    }

    if (pool_bytes > total_tensor_bytes) {
        std::fprintf(stderr, "expert pool exceeds total tensor bytes\n");
        gguf_free(ctx);
        return 8;
    }
    const uint64_t non_expert_tensor_bytes = total_tensor_bytes - pool_bytes;
    const bool uniform = min_bytes_per_expert == max_bytes_per_expert;
    const std::string uniform_json = uniform ? std::to_string(min_bytes_per_expert) : "null";

    std::fputs("{\"schema\":\"", stdout);
    std::fputs(k_schema, stdout);
    std::fputs("\",\"llama_source_commit\":\"", stdout);
    std::fputs(k_llama_commit, stdout);
    std::printf("\",\"tensor_count\":%" PRId64 ",\"total_tensor_bytes\":%" PRIu64,
        n_tensors, total_tensor_bytes);
    std::printf(",\"layers\":%d,\"experts_per_layer\":%d,\"expert_pool_bytes\":%" PRIu64,
        k_layers, k_experts, pool_bytes);
    std::printf(",\"non_expert_tensor_bytes\":%" PRIu64, non_expert_tensor_bytes);
    std::printf(",\"uniform_expert_bytes\":%s", uniform_json.c_str());
    std::printf(",\"min_expert_bytes\":%" PRIu64 ",\"max_expert_bytes\":%" PRIu64, min_bytes_per_expert, max_bytes_per_expert);
    std::fputs(",\"layer_storage\":[", stdout);

    for (size_t i = 0; i < layers.size(); ++i) {
        if (i != 0) std::putchar(',');
        const layer_info & layer = layers[i];
        std::printf("{\"layer\":%d,\"bytes\":%" PRIu64 ",\"bytes_per_expert\":%" PRIu64 ",",
            layer.layer, layer.bytes, layer.bytes_per_expert);
        print_tensor("gate", layer.gate);
        std::putchar(',');
        print_tensor("up", layer.up);
        std::putchar(',');
        print_tensor("down", layer.down);
        std::putchar('}');
    }
    std::fputs("]}\n", stdout);

    gguf_free(ctx);
    return 0;
}
