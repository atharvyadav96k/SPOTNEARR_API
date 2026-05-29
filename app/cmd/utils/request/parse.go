package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type contextKey string

const bodyContextKey contextKey = "parsed_request_body"

type Body struct {
	value any
}

func (b Body) ToString() string {
	if b.value == nil {
		return ""
	}
	if str, ok := b.value.(string); ok {
		return str
	}
	return strings.TrimSpace(fmt.Sprintf("%v", b.value))
}

func (b Body) ToInt() int {
	if b.value == nil {
		return 0
	}
	switch v := b.value.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		i, _ := strconv.Atoi(v)
		return i
	default:
		return 0
	}
}

func (b Body) ToFloat64() float64 {
	if b.value == nil {
		return 0
	}

	switch v := b.value.(type) {
	case float64:
		return v

	case float32:
		return float64(v)

	case int:
		return float64(v)

	case int64:
		return float64(v)

	case int32:
		return float64(v)

	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f

	default:
		return 0
	}
}

func GetVal(r *http.Request, key string) *Body {
	rawMap, err := parseBodyToMap(r)
	if err != nil {
		return nil
	}

	rawVal, exists := rawMap[key]
	if !exists {
		return nil
	}
	return &Body{value: rawVal}
}

func parseBodyToMap(r *http.Request) (map[string]any, error) {
	if cachedData, ok := r.Context().Value(bodyContextKey).(map[string]any); ok {
		return cachedData, nil
	}

	if r.Body == nil {
		return make(map[string]any), nil
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	_ = r.Body.Close()

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var data map[string]any
	if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, err
		}
	} else {
		data = make(map[string]any)
	}

	ctx := context.WithValue(r.Context(), bodyContextKey, data)
	*r = *r.WithContext(ctx)

	return data, nil
}
