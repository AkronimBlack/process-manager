package json

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	IsGreater  = "is_greater"
	IsLower    = "is_lower"
	IsEqual    = "is_equal"
	HttpAction = "http"

	comparingKey = "comparing"
	compareToKey = "compare_to"

	result = "result"
)

func IsPlaceholder(value string) bool {
	if strings.Contains(value, "{{") && strings.Contains(value, "}}") {
		return true
	}
	return false
}

func CleanPlaceHolder(value string) string {
	return strings.TrimSpace(strings.TrimRight(strings.TrimLeft(value, "{{"), "}}"))
}

func IsGreaterHandler(ctx context.Context, action Action, session *Session) {
	resultVariable := action.Args.GetString(result, fmt.Sprintf("%s.result", action.ActionType))

	comparing := action.Args.GetString(comparingKey)
	comparingValue := session.PlaceholderOrIntValue(comparing)

	compareTo := action.Args.GetString(compareToKey)
	compareToValue := session.PlaceholderOrIntValue(compareTo)

	session.Set(resultVariable, comparingValue > compareToValue)
}

func IsLowerHandler(ctx context.Context, action Action, session *Session) {
	resultVariable := action.Args.GetString(result, fmt.Sprintf("%s.result", action.ActionType))
	comparing := action.Args.Get(comparingKey)
	comparingValue := session.PlaceholderOrIntValue(comparing)

	compareTo := action.Args.Get(compareToKey)
	compareToValue := session.PlaceholderOrIntValue(compareTo)
	session.Set(resultVariable, comparingValue < compareToValue)
}

func IsEqualHandler(ctx context.Context, action Action, session *Session) {
	resultVariable := action.Args.GetString(result, fmt.Sprintf("%s.result", action.ActionType))

	comparing := action.Args.Get(comparingKey)
	comparingValue := session.PlaceholderOrIntValue(comparing)

	compareTo := action.Args.Get(compareToKey)
	compareToValue := session.PlaceholderOrIntValue(compareTo)

	session.Set(resultVariable, comparingValue == compareToValue)
}

func HttpHandler(ctx context.Context, action Action, session *Session) {
	targetUrl := action.Args.GetString(Url)
	resultVariable := action.Args.GetString(result, fmt.Sprintf("%s.result", action.ActionType))
	method := strings.ToUpper(action.Args.GetString(Method))
	timeout := action.Args.GetInt(Timeout)

	_, err := url.ParseRequestURI(targetUrl)
	if err != nil {
		session.Set(resultVariable, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	var req *http.Request
	req, err = http.NewRequest(method, targetUrl, nil)
	client := &http.Client{}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
	defer cancel()
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		session.Set(resultVariable, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		session.Set(resultVariable, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	session.Set(resultVariable, map[string]interface{}{
		"status":      resp.Status,
		"status_code": resp.StatusCode,
		"response":    string(body),
	})
}
