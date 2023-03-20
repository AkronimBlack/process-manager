package json

import (
	"github.com/AkronimBlack/file-executor/shared"
	"testing"
)

func TestValidateIsGrater(t *testing.T) {
	action := Action{
		ActionType: "is_greater",
		Args: map[string]interface{}{
			comparingKey: "10",
			compareToKey: "12",
		},
		OnSuccess: "test_1",
		OnFailure: "test_2",
	}
	errors := ValidateTwoValueOperators(action)
	if !errors.IsValid() {
		t.Errorf("failed falidation of is greater \n %s", shared.ToJsonPrettyString(errors))
	}
	t.Log("validation errors empty")
	action = Action{
		ActionType: "is_greater",
		Args: map[string]interface{}{
			compareToKey: "12",
		},
		OnSuccess: "test_1",
		OnFailure: "test_2",
	}
	errors = ValidateTwoValueOperators(action)
	if errors.IsValid() {
		t.Error("reports valid on missing argument")
	}
	t.Logf("validation errors \n%s", shared.ToJsonPrettyString(errors))
}

func TestValidateIsGraterReturnsErrorOnInvalidAction(t *testing.T) {
	actions := []Action{
		{
			ActionType: "is_greater",
			Args: map[string]interface{}{
				compareToKey: "12",
			},
			OnSuccess: "test_1",
			OnFailure: "test_2",
		},
		{
			ActionType: "is_greater",
			Args: map[string]interface{}{
				comparingKey: "12",
			},
			OnSuccess: "test_1",
			OnFailure: "test_2",
		},
	}
	for _, action := range actions {
		errors := ValidateTwoValueOperators(action)
		if errors.IsValid() {
			t.Error("reports valid on missing argument")
		}
		t.Logf("validation errors \n%s", shared.ToJsonPrettyString(errors))
	}
}

func TestValidateHttpAction(t *testing.T) {
	action := Action{
		ActionType: HttpAction,
		Args: map[string]interface{}{
			Url:     "value",
			Method:  "value",
			Timeout: 100,
			Headers: map[string]interface{}{
				"test": "test",
			},
			Payload: map[string]interface{}{
				"test_1": "test_1",
				"test_2": "test_2",
				"test_3": "test_3",
			},
			result: "http_action_result",
		},
		OnSuccess: "test_1",
		OnFailure: "test_2",
	}
	errors := ValidateHttpAction(action)
	if !errors.IsValid() {
		t.Error("reports invalid on valid args")
	}
	t.Logf("validation errors \n%s", shared.ToJsonPrettyString(errors))
}

func TestValidateHttpActionFail(t *testing.T) {
	actions := []Action{
		{
			ActionType: HttpAction,
			Args: map[string]interface{}{
				Url:     "value",
				Method:  "value",
				Timeout: "value",
				Headers: map[string]interface{}{
					"test": "test",
				},
			},
			OnSuccess: "test_1",
			OnFailure: "test_2",
		},
		{
			ActionType: HttpAction,
			Args: map[string]interface{}{
				Url:     "value",
				Method:  "value",
				Timeout: "value",
				Payload: map[string]interface{}{
					"test_1": "test_1",
					"test_2": "test_2",
					"test_3": "test_3",
				},
			},
			OnSuccess: "test_1",
			OnFailure: "test_2",
		},
	}
	for _, action := range actions {
		errors := ValidateHttpAction(action)
		if errors.IsValid() {
			t.Error("reports valid on invalid args")
		}
		t.Logf("validation errors \n%s", shared.ToJsonPrettyString(errors))
	}
}
