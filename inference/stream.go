package inference

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"
)

// NewStreamingLocalHTTP enables bounded SSE collection with measured first-content latency.
// The returned result is delivered only after a complete, validated stream.
func NewStreamingLocalHTTP(id, model, endpoint string, features []string) (*LocalHTTP, error) {
	b, e := NewLocalHTTP(id, model, endpoint, features)
	if e == nil {
		b.streaming = true
	}
	return b, e
}
func (b *LocalHTTP) Close() { b.client.CloseIdleConnections() }
func usageValid(u Usage, max int) bool {
	return (u.PromptTokens == nil || *u.PromptTokens >= 0) && (u.OutputTokens == nil || *u.OutputTokens >= 0 && *u.OutputTokens <= max)
}
func (b *LocalHTTP) readStream(body io.Reader, r Request, start time.Time) (Result, error) {
	out := Result{BackendID: b.name, ModelID: b.model, Status: "failed"}
	limit := int64(r.MaxOutputTokens*256 + 65536)
	reader := &io.LimitedReader{R: body, N: limit + 1}
	scan := bufio.NewScanner(reader)
	scan.Buffer(make([]byte, 4096), r.MaxOutputTokens*16+65536)
	seenUsage := false
	finished := false
	done := false
	data := ""
	var first, last time.Time
	event := func() error {
		if data == "" {
			return nil
		}
		payload := strings.TrimSuffix(data, "\n")
		data = ""
		if payload == "[DONE]" {
			if done {
				return fmt.Errorf("duplicate completion")
			}
			if !finished {
				return fmt.Errorf("stream ended before stop")
			}
			done = true
			return nil
		}
		if done {
			return fmt.Errorf("data after done")
		}
		if !utf8.ValidString(payload) {
			return fmt.Errorf("invalid UTF-8 stream")
		}
		var chunk struct {
			Model   string          `json:"model"`
			Error   json.RawMessage `json:"error"`
			Choices []struct {
				Index int `json:"index"`
				Delta struct {
					Content string          `json:"content"`
					Tools   json.RawMessage `json:"tool_calls"`
				} `json:"delta"`
				Finish *string `json:"finish_reason"`
			} `json:"choices"`
			Usage *struct {
				Prompt *int `json:"prompt_tokens"`
				Output *int `json:"completion_tokens"`
			} `json:"usage"`
			Timings *struct {
				PromptMS           float64 `json:"prompt_ms"`
				PromptPerSecond    float64 `json:"prompt_per_second"`
				PredictedPerSecond float64 `json:"predicted_per_second"`
			} `json:"timings"`
		}
		if e := json.Unmarshal([]byte(payload), &chunk); e != nil {
			return e
		}
		if chunk.Model != b.model || len(chunk.Error) > 0 && string(chunk.Error) != "null" || len(chunk.Choices) > 1 {
			return fmt.Errorf("invalid stream model/error/choices")
		}
		if chunk.Usage != nil {
			if seenUsage {
				return fmt.Errorf("duplicate usage event")
			}
			seenUsage = true
			out.Usage = Usage{chunk.Usage.Prompt, chunk.Usage.Output}
			if !usageValid(out.Usage, r.MaxOutputTokens) {
				return fmt.Errorf("stream token budget exceeded")
			}
		}
		if chunk.Timings != nil {
			if chunk.Timings.PromptMS < 0 || chunk.Timings.PromptPerSecond < 0 || chunk.Timings.PredictedPerSecond < 0 {
				return fmt.Errorf("invalid stream timings")
			}
			if chunk.Timings.PromptMS > 0 {
				v := chunk.Timings.PromptMS
				out.Telemetry.PromptEvalMS = &v
			}
			if chunk.Timings.PromptPerSecond > 0 {
				v := chunk.Timings.PromptPerSecond
				out.Telemetry.PromptTokensPerSecond = &v
			}
			if chunk.Timings.PredictedPerSecond > 0 {
				v := chunk.Timings.PredictedPerSecond
				out.Telemetry.ServerGenerationTokensPerSecond = &v
			}
		}
		if len(chunk.Choices) == 1 {
			c := chunk.Choices[0]
			if c.Index != 0 || len(c.Delta.Tools) > 0 && string(c.Delta.Tools) != "null" {
				return fmt.Errorf("unsupported stream choice/tool call")
			}
			if finished && (c.Delta.Content != "" || c.Finish != nil) {
				return fmt.Errorf("content after finish")
			}
			if c.Delta.Content != "" {
				last = time.Now()
				if first.IsZero() {
					first = last
				}
				if len(out.Text)+len(c.Delta.Content) > r.MaxOutputTokens*16 {
					return fmt.Errorf("stream output byte limit")
				}
				out.Text += c.Delta.Content
			}
			if c.Finish != nil {
				if *c.Finish != "stop" {
					return fmt.Errorf("incomplete stream: %s", *c.Finish)
				}
				finished = true
			}
		}
		return nil
	}
	for scan.Scan() {
		line := scan.Text()
		if reader.N <= 0 {
			return out, fmt.Errorf("stream total byte limit")
		}
		if line == "" {
			if e := event(); e != nil {
				return out, e
			}

			continue
		}
		if strings.HasPrefix(line, "data:") {
			part := strings.TrimPrefix(line, "data:")
			part = strings.TrimPrefix(part, " ")
			if len(data)+len(part) > r.MaxOutputTokens*16+65536 {
				return out, fmt.Errorf("stream event byte limit")
			}
			data += part + "\n"
		} else if !strings.HasPrefix(line, ":") && !strings.HasPrefix(line, "event:") && !strings.HasPrefix(line, "id:") && !strings.HasPrefix(line, "retry:") {
			return out, fmt.Errorf("invalid SSE field")
		}
	}
	if e := scan.Err(); e != nil {
		return out, e
	}
	if !done {
		return out, fmt.Errorf("truncated stream")
	}
	if !first.IsZero() {
		ttft := first.Sub(start).Milliseconds()
		generation := last.Sub(first).Milliseconds()
		out.Telemetry.TimeToFirstTokenMS = &ttft
		out.Telemetry.GenerationMS = &generation
	}
	out.Status = "completed"
	out.Termination = "stop"
	return out, nil
}
