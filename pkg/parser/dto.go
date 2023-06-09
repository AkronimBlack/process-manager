package parser

func NewSessionDto(session Session) SessionDto {
	if session == nil {
		return SessionDto{}
	}
	return SessionDto{
		Uuid:                    session.Uuid(),
		Values:                  session.Values(),
		ExecutedActions:         NewExecutedActionsDto(session.ExecutedActions()),
		InputData:               session.InputData(),
		OnFinishWebhook:         NewOnFinishWebhookDto(session.OnFinishWebhook()),
		OnFinishWebhookResponse: session.OnFinishWebhookResponse(),
		Tasks:                   NewTasksDto(session.Tasks()),
	}
}

type SessionDto struct {
	Uuid                    string                 `json:"uuid"`
	Values                  map[string]interface{} `json:"values"`
	ExecutedActions         []ExecutedActionDto    `json:"executed_actions"`
	InputData               map[string]interface{} `json:"input_data"`
	OnFinishWebhook         *OnFinishWebhook       `json:"on_finish_webhook"`
	OnFinishWebhookResponse map[string]interface{} `json:"on_finish_webhook_response"`
	Tasks                   []TaskDto              `json:"tasks"`
}

func NewOnFinishWebhookDto(onFinishWebhook Webhook) *OnFinishWebhook {
	if onFinishWebhook == nil {
		return nil
	}
	return &OnFinishWebhook{
		Url: onFinishWebhook.Url(),
	}
}

type OnFinishWebhook struct {
	Url string `json:"url"`
}

func NewExecutedActionsDto(executedActions []ExecutedAction) []ExecutedActionDto {
	values := make([]ExecutedActionDto, len(executedActions))
	for i, v := range executedActions {
		values[i] = NewExecutedActionDto(v)
	}
	return values
}

func NewExecutedActionDto(executedAction ExecutedAction) ExecutedActionDto {
	return ExecutedActionDto{
		ActionType: executedAction.Type(),
		Args:       executedAction.Arguments(),
		OnSuccess:  executedAction.OnSuccess(),
		OnFailure:  executedAction.OnFailure(),
		Params:     executedAction.Parameters(),
	}
}

type ExecutedActionDto struct {
	ActionType string `json:"type"`
	Args       Args   `json:"args"`
	OnSuccess  string `json:"on_success"`
	OnFailure  string `json:"on_failure"`
	Params     map[string]interface{}
}

func NewTasksDto(tasks []Task) []TaskDto {
	tasksDto := make([]TaskDto, len(tasks))
	for i, v := range tasks {
		tasksDto[i] = NewTaskDto(v)
	}
	return tasksDto
}

func NewTaskDto(task Task) TaskDto {
	return TaskDto{
		ID:         task.ID(),
		Name:       task.Name(),
		Next:       task.Next(),
		Parameters: task.Parameters(),
	}
}

type TaskDto struct {
	ID         string                 `json:"ID"`
	Name       string                 `json:"name"`
	Next       string                 `json:"next"`
	Parameters map[string]interface{} `json:"parameters"`
}
