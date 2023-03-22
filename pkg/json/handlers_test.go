package json

import (
	"context"
	"testing"
)

func TestIsPlaceholder(t *testing.T) {
	value := "{{session.test}}"
	value2 := "session.test"
	if !IsPlaceholder(value) {
		t.Errorf("%s not recognized as placeholder", value)
	}
	if IsPlaceholder(value2) {
		t.Errorf("%s recognized as placeholder", value2)
	}
}

func TestCleanPlaceHolder(t *testing.T) {
	value := "{{  session.test    }}"
	cleaned := CleanPlaceHolder(value)
	if cleaned != "session.test" {
		t.Errorf("placeholder %s clean error, result %s", value, cleaned)
	}
}

func TestIsGreaterHandler(t *testing.T) {
	actions := []Action{
		{
			ActionType: IsGreater,
			Args: map[string]interface{}{
				comparingKey: 10,
				compareToKey: 11,
				result:       "test_result",
			},
			OnSuccess: "",
			OnFailure: "",
		},
	}

	for _, action := range actions {
		session := &Session{
			values:          map[string]interface{}{},
			executedActions: []Action{},
		}
		err := IsGreaterHandler(context.Background(), action, session)
		if err != nil {
			t.Errorf("error while executing handler: %s", err.Error())
			return
		}
		if session.ValueOf("test_result").(bool) {
			t.Errorf("wrong evaluation of 10>11, %v", session.ValueOf("test_result").(bool))
			return
		}
		t.Logf("correct evaluation of 10>11, %v", session.ValueOf("test_result").(bool))
	}
}

func TestIsLowerHandler(t *testing.T) {
	actions := []Action{
		{
			ActionType: IsLower,
			Args: map[string]interface{}{
				comparingKey: 10,
				compareToKey: 11,
				result:       "test_result",
			},
			OnSuccess: "",
			OnFailure: "",
		},
	}

	for _, action := range actions {
		session := &Session{
			values:          map[string]interface{}{},
			executedActions: []Action{},
		}
		err := IsLowerHandler(context.Background(), action, session)
		if err != nil {
			t.Errorf("error while executing handler: %s", err.Error())
			return
		}
		if !session.ValueOf("test_result").(bool) {
			t.Errorf("wrong evaluation of 10<11, %v", session.ValueOf("test_result").(bool))
			return
		}
		t.Logf("correct evaluation of 10<11, %v", session.ValueOf("test_result").(bool))
	}
}

func TestIsEqualHandlerHandler(t *testing.T) {
	actions := []Action{
		{
			ActionType: IsEqual,
			Args: map[string]interface{}{
				comparingKey: 10,
				compareToKey: 10,
				result:       "test_result",
			},
			OnSuccess: "",
			OnFailure: "",
		},
	}

	for _, action := range actions {
		session := &Session{
			values:          map[string]interface{}{},
			executedActions: []Action{},
		}
		err := IsEqualHandler(context.Background(), action, session)
		if err != nil {
			t.Errorf("error while executing handler: %s", err.Error())
			return
		}
		if !session.ValueOf("test_result").(bool) {
			t.Errorf("wrong evaluation of 10==11, %v", session.ValueOf("test_result").(bool))
			return
		}
		t.Logf("correct evaluation of 10==11, %v", session.ValueOf("test_result").(bool))
	}
}

func TestHttpHandler(t *testing.T) {
	action := Action{
		ActionType: HttpAction,
		Args: map[string]interface{}{
			"url":     "https://api.myip.com",
			"method":  "get",
			"timeout": 400,
			"headers": map[string]interface{}{},
			"payload": map[string]interface{}{},
			result:    "http_action_result",
		},
		OnSuccess: "test_1",
		OnFailure: "test_2",
	}
	session := &Session{
		values:          map[string]interface{}{},
		executedActions: []Action{},
	}
	err := HttpHandler(context.Background(), action, session)
	if err != nil {
		t.Errorf("error while executing handler: %s", err.Error())
		return
	}
	httpActionError := session.StringValueOf("http_action_result.error", "")
	if httpActionError != "" {
		t.Errorf("http action failed with error %s", httpActionError)
		return
	}
	httpResult := session.ValueOf("http_action_result")
	if httpResult == nil {
		t.Error("http action result empty")
	}
}
