package json

import (
	"github.com/AkronimBlack/file-executor/shared"
	"testing"
)

func TestParser_LoadFile(t *testing.T) {
	location := "test.json"
	parser := Parser{}
	err := parser.LoadFile(location)
	if err != nil {
		t.Error(err)
	}
	if len(parser.Actions()) == 0 {
		t.Errorf("failed loading %s, actions are 0", location)
	}
}

func TestParser_LoadFileFailsToLoadNonJsonFile(t *testing.T) {
	location := "test.xml"
	parser := Parser{}
	err := parser.LoadFile(location)
	if err == nil {
		t.Error("parser did not throw an error while loading an .xml file")
	}
}

func TestParser_LoadFileFailsToLoadNonExistingFile(t *testing.T) {
	location := "non_existing_json_file.json"
	parser := Parser{}
	err := parser.LoadFile(location)
	if err == nil {
		t.Error("parser did not throw an error while loading an non_existing_json_file.json file")
	}
}

func TestParser_ValidateAction(t *testing.T) {
	parser := Parser{}
	validationErrors := parser.ValidateAction(Action{
		ActionType: "sum",
		Args:       map[string]interface{}{},
		OnSuccess:  "test_id",
		OnFailure:  "test_id",
	})
	if len(validationErrors) != 0 {
		t.Errorf("found validation errors on valid action\n%s", shared.ToJsonPrettyString(validationErrors))
	}
}

func TestParser_ValidateActionReturnsErrorsOnInvalidAction(t *testing.T) {
	actions := []Action{
		{
			ActionType: "",
			Args:       map[string]interface{}{},
			OnSuccess:  "test_id",
			OnFailure:  "test_id",
		},
		{
			ActionType: "sum",
			Args:       map[string]interface{}{},
			OnSuccess:  "",
			OnFailure:  "test_id",
		},
		{
			ActionType: "sum",
			Args:       map[string]interface{}{},
			OnSuccess:  "test_id",
			OnFailure:  "",
		},
	}
	parser := Parser{}
	for _, action := range actions {
		validationErrors := parser.ValidateAction(action)
		if len(validationErrors) == 0 {
			t.Error("did not find validation errors on invalid action")
		}
		t.Logf("found validation errors \n%s", shared.ToJsonPrettyString(validationErrors))
	}
}

func TestParser_validate(t *testing.T) {
	actions := map[string]Action{
		"test_id_1": {
			ActionType: "sum",
			Args:       map[string]interface{}{},
			OnSuccess:  "test_id",
			OnFailure:  "test_id",
		},
		"test_id_2": {
			ActionType: "sum",
			Args:       map[string]interface{}{},
			OnSuccess:  "test_id",
			OnFailure:  "test_id",
		},
		"test_id_3": {
			ActionType: "sum",
			Args:       map[string]interface{}{},
			OnSuccess:  "test_id",
			OnFailure:  "test_id",
		},
	}
	parser := Parser{}
	validationErrors := parser.validate(actions)
	if !validationErrors.IsValid() {
		t.Errorf("found validation errors on valid action\n%s", shared.ToJsonPrettyString(validationErrors))
	}
}

func TestParser_validateReturnsErrorsOnInvalidAction(t *testing.T) {
	actions := map[string]Action{
		"test_id_1": {
			ActionType: "sum",
			Args:       map[string]interface{}{},
			OnSuccess:  "test_id",
			OnFailure:  "test_id",
		},
		"test_id_2": {
			ActionType: "",
			Args:       map[string]interface{}{},
			OnSuccess:  "test_id",
			OnFailure:  "test_id",
		},
		"test_id_3": {
			ActionType: "sum",
			Args:       map[string]interface{}{},
			OnSuccess:  "test_id",
			OnFailure:  "test_id",
		},
	}
	parser := Parser{}
	validationErrors := parser.validate(actions)
	if validationErrors.IsValid() {
		t.Error("did not find validation errors on invalid action")
	}
	t.Logf("found validation errors \n%s", shared.ToJsonPrettyString(validationErrors))
}

func TestParser_Validate(t *testing.T) {
	actions := map[string]Action{
		"test_id_1": {
			ActionType: "sum",
			Args:       map[string]interface{}{},
			OnSuccess:  "test_id",
			OnFailure:  "test_id",
		},
		"test_id_2": {
			ActionType: "sum",
			Args:       map[string]interface{}{},
			OnSuccess:  "test_id",
			OnFailure:  "test_id",
		},
		"test_id_3": {
			ActionType: "sum",
			Args:       map[string]interface{}{},
			OnSuccess:  "test_id",
			OnFailure:  "test_id",
		},
	}
	parser := Parser{}
	parser.SetActions(actions)
	validationErrors := parser.Validate()
	if !validationErrors.IsValid() {
		t.Errorf("found validation errors on valid action\n%s", shared.ToJsonPrettyString(validationErrors))
	}
}

func TestParser_ValidateReturnsErrorsFromCustomAndDefaultValidation(t *testing.T) {
	actions := map[string]Action{
		"test_id_1": {
			ActionType: "is_greater",
			Args:       map[string]interface{}{},
			OnSuccess:  "",
			OnFailure:  "test_id",
		},
	}
	parser := Parser{
		validators: map[string]Validator{},
	}
	parser.AddValidator("is_greater", ValidateTwoValueOperators)
	parser.SetActions(actions)
	validationErrors := parser.Validate()
	if validationErrors.IsValid() {
		t.Error("did not find errors on an invalid action")
	}
	t.Logf("validation errors\n%s", shared.ToJsonPrettyString(validationErrors))
}
