package parser

import "context"

// Handler interface. Expects id of next action to execute. Returning empty string finished the process execution
type Handler func(ctx context.Context, action *Action, session Session) string

type Args map[string]interface{}

func (a Args) Get(key string) interface{} {
	return a[key]
}

type Session interface {
	Uuid() string
	Values() map[string]interface{}
	ExecutedActions() []ExecutedAction
	InputData() map[string]interface{}
	SetInputData(inputData map[string]interface{})
	Set(key string, value interface{})
	AddExecutedAction(action ExecutedAction)
	OnFinishWebhook() Webhook
	OnFinishWebhookResponse() map[string]interface{}
	SetOnFinishWebhook(onFinishWebhook Webhook)
	SetOnFinishWebhookResponse(onFinishWebhookResponse map[string]interface{})
	PlaceholderOrStringValue(value string) string
	PlaceholderOrIntValue(value interface{}) int64
	ValueOf(key string) interface{}
	StringValueOf(key string, defaultValue string) string
	IntValueOf(key string, defaultValue int64) int64
	Tasks() []Task
	Task(id string) Task
	AddTask(task Task)
	UpdateData(parameters map[string]interface{})
}

type ExecutedAction interface {
	Type() string
	Arguments() Args
	OnSuccess() string
	OnFailure() string
	Parameters() map[string]interface{}
}

type Webhook interface {
	Url() string
}

type Task interface {
	ID() string
	Name() string
	Next() string
	Parameters() map[string]interface{}
	Execute(parameters map[string]interface{})
	Session() Session
}
