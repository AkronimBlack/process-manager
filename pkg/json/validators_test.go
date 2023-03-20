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
