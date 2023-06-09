package parser

import (
	"context"
	"github.com/AkronimBlack/process-manager/shared"
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
	validationErrors := parser.ValidateAction(&Action{
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
	actions := []*Action{
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
	actions := map[string]*Action{
		StartNode: {
			ActionType: StartNode,
			OnSuccess:  "test_id_1",
		},
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
	actions := map[string]*Action{
		StartNode: {
			ActionType: StartNode,
			OnSuccess:  "test_id_1",
		},
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
	actions := map[string]*Action{
		StartNode: {
			ActionType: StartNode,
			OnSuccess:  "test_id_1",
		},
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

func TestParser_Execute(t *testing.T) {
	actions := map[string]*Action{
		StartNode: {
			ActionType: StartNode,
			OnSuccess:  "test_id_1",
		},
		"test_id_1": {
			ActionType: IsGreater,
			Args: map[string]interface{}{
				comparingKey: 10,
				compareToKey: 11,
				result:       "test_result_1",
			},
			OnSuccess: "test_id_2",
			OnFailure: "test_id_3",
		},
		"test_id_2": {
			ActionType: IsLower,
			Args: map[string]interface{}{
				comparingKey: "{{$session.values.test_result_1}}",
				compareToKey: 11,
				result:       "test_result_2",
			},
			OnSuccess: "test_1",
			OnFailure: "test_2",
		},
		"test_id_3": {
			ActionType: IsEqual,
			Args: map[string]interface{}{
				comparingKey: 10,
				compareToKey: 10,
				result:       "test_result_3",
			},
			OnSuccess: "test_1",
			OnFailure: "test_2",
		},
	}
	parser := Parser{
		handlers: map[string]Handler{
			IsGreater: IsGreaterHandler,
			IsLower:   IsLowerHandler,
			IsEqual:   IsEqualHandler,
		},
		sessions: make([]Session, 0),
	}
	parser.SetActions(actions)
	validationErrors := parser.Validate()
	if !validationErrors.IsValid() {
		t.Errorf("found validation errors on valid action\n%s", shared.ToJsonPrettyString(validationErrors))
	}
	data := map[string]interface{}{
		"data": "data",
	}
	sessionUuid := parser.Execute(context.Background(), data, nil)
	if sessionUuid == "" {
		t.Error("session_uuid is empty")
		return
	}
}
