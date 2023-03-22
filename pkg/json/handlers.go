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

type ResultArgs struct {
	Result string `json:"result"`
}

func (a ResultArgs) ResultVariable(actionType string) string {
	if a.Result == "" {
		return fmt.Sprintf("%s.result", actionType)
	}
	return a.Result
}

type OperatorArgs struct {
	ResultArgs
	Comparing int `json:"comparing"`
	CompareTo int `json:"compare_to"`
}

func IsGreaterHandler(ctx context.Context, action Action, session *Session) error {
	operatorArgs := OperatorArgs{}
	err := action.Args.Bind(&operatorArgs)
	if err != nil {
		return err
	}
	session.Set(
		operatorArgs.ResultVariable(action.ActionType),
		session.PlaceholderOrIntValue(operatorArgs.Comparing) > session.PlaceholderOrIntValue(operatorArgs.CompareTo),
	)
	return nil
}

func IsLowerHandler(ctx context.Context, action Action, session *Session) error {
	operatorArgs := OperatorArgs{}
	err := action.Args.Bind(&operatorArgs)
	if err != nil {
		return err
	}
	session.Set(
		operatorArgs.ResultVariable(action.ActionType),
		session.PlaceholderOrIntValue(operatorArgs.Comparing) < session.PlaceholderOrIntValue(operatorArgs.CompareTo),
	)
	return nil
}

func IsEqualHandler(ctx context.Context, action Action, session *Session) error {
	operatorArgs := OperatorArgs{}
	err := action.Args.Bind(&operatorArgs)
	if err != nil {
		return err
	}
	session.Set(
		operatorArgs.ResultVariable(action.ActionType),
		session.PlaceholderOrIntValue(operatorArgs.Comparing) == session.PlaceholderOrIntValue(operatorArgs.CompareTo),
	)
	return nil
}

type HttpHandlerArgs struct {
	ResultArgs
	Url        string `json:"url"`
	HttpMethod string `json:"method"`
	Timeout    int    `json:"timeout"`
}

func (h HttpHandlerArgs) Method() string {
	return strings.ToUpper(h.HttpMethod)
}

func HttpHandler(ctx context.Context, action Action, session *Session) error {
	httpArgs := HttpHandlerArgs{}
	err := action.Args.Bind(&httpArgs)
	if err != nil {
		return err
	}
	_, err = url.ParseRequestURI(httpArgs.Url)
	if err != nil {
		session.Set(httpArgs.ResultVariable(action.ActionType), map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	var req *http.Request
	req, err = http.NewRequest(httpArgs.Method(), httpArgs.Url, nil)
	client := &http.Client{}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(httpArgs.Timeout)*time.Millisecond)
	defer cancel()
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		session.Set(httpArgs.ResultVariable(action.ActionType), map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		session.Set(httpArgs.ResultVariable(action.ActionType), map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	session.Set(httpArgs.ResultVariable(action.ActionType), map[string]interface{}{
		"status":      resp.Status,
		"status_code": resp.StatusCode,
		"response":    string(body),
	})
	return nil
}
