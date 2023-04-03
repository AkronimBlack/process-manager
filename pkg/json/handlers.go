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

func (a ResultArgs) ResultVariableAsError(actionType string) string {
	if a.Result == "" {
		return fmt.Sprintf("%s.result_error", actionType)
	}
	return fmt.Sprintf("%s_error", a.Result)
}

type OperatorArgs struct {
	ResultArgs
	Comparing string `json:"comparing"`
	CompareTo string `json:"compare_to"`
}

func AddActionError(session *Session, variable string, err error) {
	if session == nil || variable == "" {
		return
	}
	session.Set(variable, err.Error())
}

func IsGreaterHandler(ctx context.Context, action *Action, session *Session) string {
	operatorArgs := OperatorArgs{}
	err := action.Args.Bind(&operatorArgs)
	if err != nil {
		AddActionError(session, operatorArgs.ResultVariableAsError(action.ActionType), err)
		return action.OnFailure
	}
	comparing := session.PlaceholderOrIntValue(operatorArgs.Comparing)
	compareTo := session.PlaceholderOrIntValue(operatorArgs.CompareTo)
	session.Set(
		operatorArgs.ResultVariable(action.ActionType),
		comparing > compareTo,
	)
	session.executedActions = append(session.executedActions, &ExecutedAction{
		Action: *action,
		Params: map[string]interface{}{
			"comparing":  comparing,
			"compare_to": compareTo,
		},
	})
	return action.OnSuccess
}

func IsLowerHandler(ctx context.Context, action *Action, session *Session) string {
	operatorArgs := OperatorArgs{}
	err := action.Args.Bind(&operatorArgs)
	if err != nil {
		AddActionError(session, operatorArgs.ResultVariableAsError(action.ActionType), err)
		return action.OnFailure
	}
	session.Set(
		operatorArgs.ResultVariable(action.ActionType),
		session.PlaceholderOrIntValue(operatorArgs.Comparing) < session.PlaceholderOrIntValue(operatorArgs.CompareTo),
	)
	return action.OnSuccess
}

func IsEqualHandler(ctx context.Context, action *Action, session *Session) string {
	operatorArgs := OperatorArgs{}
	err := action.Args.Bind(&operatorArgs)
	if err != nil {
		AddActionError(session, operatorArgs.ResultVariableAsError(action.ActionType), err)
		return action.OnFailure
	}
	session.Set(
		operatorArgs.ResultVariable(action.ActionType),
		session.PlaceholderOrIntValue(operatorArgs.Comparing) == session.PlaceholderOrIntValue(operatorArgs.CompareTo),
	)
	return action.OnSuccess
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

func HttpHandler(ctx context.Context, action *Action, session *Session) string {
	httpArgs := HttpHandlerArgs{}
	err := action.Args.Bind(&httpArgs)
	if err != nil {
		AddActionError(session, httpArgs.ResultVariableAsError(action.ActionType), err)
		return action.OnFailure
	}
	_, err = url.ParseRequestURI(httpArgs.Url)
	if err != nil {
		AddActionError(session, httpArgs.ResultVariableAsError(action.ActionType), err)
		return action.OnFailure
	}

	var req *http.Request
	req, err = http.NewRequest(httpArgs.Method(), httpArgs.Url, nil)
	client := &http.Client{}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(httpArgs.Timeout)*time.Millisecond)
	defer cancel()
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		AddActionError(session, httpArgs.ResultVariableAsError(action.ActionType), err)
		return action.OnFailure
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		AddActionError(session, httpArgs.ResultVariableAsError(action.ActionType), err)
		return action.OnFailure
	}
	session.Set(httpArgs.ResultVariable(action.ActionType), map[string]interface{}{
		"status":      resp.Status,
		"status_code": resp.StatusCode,
		"response":    string(body),
	})
	return action.OnSuccess
}
