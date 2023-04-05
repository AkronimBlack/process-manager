package json

func NewSessionDto(session *Session) SessionDto {
	if session == nil {
		return SessionDto{}
	}
	return SessionDto{
		Uuid:                    session.Uuid,
		Values:                  session.values,
		ExecutedActions:         NewExecutedActionsDto(session.executedActions),
		InputData:               session.inputData,
		OnFinishWebhook:         NewOnFinishWebhookDto(session.OnFinishWebhook),
		OnFinishWebhookResponse: session.OnFinishWebhookResponse,
	}
}

type SessionDto struct {
	Uuid                    string                 `json:"uuid"`
	Values                  map[string]interface{} `json:"values"`
	ExecutedActions         []ExecutedActionDto    `json:"executed_actions"`
	InputData               map[string]interface{} `json:"input_data"`
	OnFinishWebhook         *OnFinishWebhook       `json:"on_finish_webhook"`
	OnFinishWebhookResponse map[string]interface{} `json:"on_finish_webhook_response,omitempty"`
}

func NewOnFinishWebhookDto(onFinishWebhook *Webhook) *OnFinishWebhook {
	if onFinishWebhook == nil {
		return nil
	}
	return &OnFinishWebhook{
		Url: onFinishWebhook.Url,
	}
}

type OnFinishWebhook struct {
	Url string `json:"url"`
}

func NewExecutedActionsDto(executedActions []*ExecutedAction) []ExecutedActionDto {
	values := make([]ExecutedActionDto, len(executedActions))
	for i, v := range executedActions {
		values[i] = NewExecutedActionDto(v)
	}
	return values
}

func NewExecutedActionDto(executedAction *ExecutedAction) ExecutedActionDto {
	return ExecutedActionDto{
		ActionType: executedAction.ActionType,
		Args:       executedAction.Args,
		OnSuccess:  executedAction.OnSuccess,
		OnFailure:  executedAction.OnFailure,
		Params:     executedAction.Params,
	}
}

type ExecutedActionDto struct {
	ActionType string `json:"type"`
	Args       Args   `json:"args"`
	OnSuccess  string `json:"on_success"`
	OnFailure  string `json:"on_failure"`
	Params     map[string]interface{}
}
