package parser

import (
	"encoding/json"
	"github.com/AkronimBlack/process-manager/shared"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"log"
	"strconv"
	"sync"
)

func (a Args) GetString(key string, defaultValue ...string) string {
	if defaultValue == nil || len(defaultValue) == 0 {
		defaultValue = []string{""}
	}
	v, ok := a[key]
	if !ok {
		return defaultValue[0]
	}
	value, ok := v.(string)
	if !ok {
		return defaultValue[0]
	}
	return value
}

// Bind unpack args into a target struct.
func (a Args) Bind(target interface{}) error {
	return json.Unmarshal(shared.ToJsonByte(a), &target)
}

func (a Args) GetInt(key string, defaultValue ...int) int {
	if defaultValue == nil || len(defaultValue) == 0 {
		defaultValue = []int{0}
	}
	v, ok := a[key]
	if !ok {
		return defaultValue[0]
	}
	value, ok := v.(int)
	if !ok {
		return defaultValue[0]
	}
	return value
}

func (a Args) GetMap(key string, defaultValue ...map[string]interface{}) map[string]interface{} {
	if defaultValue == nil || len(defaultValue) == 0 {
		defaultValue = make([]map[string]interface{}, 1)
	}
	v, ok := a[key]
	if !ok {
		return defaultValue[0]
	}
	value, ok := v.(map[string]interface{})
	if !ok {
		return defaultValue[0]
	}
	return value
}

type Action struct {
	ActionType string `json:"type"`
	Args       Args   `json:"args"`
	OnSuccess  string `json:"on_success"`
	OnFailure  string `json:"on_failure"`
}

type executedAction struct {
	Action
	Params map[string]interface{}
}

func (e executedAction) Type() string {
	return e.ActionType
}

func (e executedAction) Arguments() Args {
	return e.Args
}

func (e executedAction) OnSuccess() string {
	return e.Action.OnSuccess
}

func (e executedAction) OnFailure() string {
	return e.Action.OnFailure
}

func (e executedAction) Parameters() map[string]interface{} {
	return e.Params
}

type webHook struct {
	url string
}

func NewWebHook(url string) Webhook {
	return &webHook{url: url}
}

func (w webHook) Url() string {
	return w.url
}

type session struct {
	uuid                    string
	values                  map[string]interface{}
	executedActions         []ExecutedAction
	inputData               map[string]interface{}
	tasks                   []Task
	onFinishWebhook         Webhook
	onFinishWebhookResponse map[string]interface{}

	lock sync.Mutex
}

func (s *session) Tasks() []Task {
	return s.tasks
}

func (s *session) SetOnFinishWebhook(onFinishWebhook Webhook) {
	s.onFinishWebhook = onFinishWebhook
}

func (s *session) SetOnFinishWebhookResponse(onFinishWebhookResponse map[string]interface{}) {
	s.onFinishWebhookResponse = onFinishWebhookResponse
}

func (s *session) Uuid() string {
	return s.uuid
}

func (s *session) Values() map[string]interface{} {
	return s.values
}

func (s *session) ExecutedActions() []ExecutedAction {
	return s.executedActions
}

func (s *session) InputData() map[string]interface{} {
	return s.inputData
}

func (s *session) OnFinishWebhook() Webhook {
	return s.onFinishWebhook
}

func (s *session) OnFinishWebhookResponse() map[string]interface{} {
	return s.onFinishWebhookResponse
}

func NewSession(data map[string]interface{}, webhook Webhook) Session {
	return &session{
		uuid:            uuid.NewString(),
		values:          make(map[string]interface{}),
		executedActions: make([]ExecutedAction, 0),
		tasks:           make([]Task, 0),
		onFinishWebhook: webhook,
		inputData:       data,
	}
}

func (s *session) AddTask(task Task) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.tasks = append(s.tasks, task)
}

func (s *session) AddExecutedAction(action ExecutedAction) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.executedActions = append(s.executedActions, action)
}

func (s *session) ValueOf(key string) interface{} {
	value := gjson.Get(shared.ToJsonString(s.values), key)
	if value.Exists() {
		return value.Value()
	}
	return nil
}

func (s *session) StringValueOf(key string, defaultValue string) string {
	value := gjson.Get(shared.ToJsonString(NewSessionDto(s)), key)
	if value.Exists() {
		return value.String()
	}
	return defaultValue
}

func (s *session) IntValueOf(key string, defaultValue int64) int64 {
	log.Println(shared.ToJsonPrettyString(NewSessionDto(s)))
	value := gjson.Get(shared.ToJsonString(NewSessionDto(s)), key)
	if value.Exists() {
		return value.Int()
	}
	return defaultValue
}

func (s *session) PlaceholderOrStringValue(value string) string {
	if IsPlaceholder(value) {
		return s.StringValueOf(CleanPlaceHolder(value), value)
	}
	return value
}

func (s *session) PlaceholderOrIntValue(value interface{}) int64 {
	switch v := value.(type) {
	case string:
		if IsPlaceholder(v) {
			return s.IntValueOf(CleanPlaceHolder(v), 0)
		}
		n, _ := strconv.ParseInt(v, 10, 64)
		return n
	case int:
		return int64(v)
	}
	return 0
}

func (s *session) Set(key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.values[key] = value
}

type task struct {
	name       string
	next       string
	parameters map[string]interface{}
}

func NewTask(name, next string, parameters map[string]interface{}) Task {
	return &task{name: name, next: next, parameters: parameters}
}

func (t task) Name() string {
	return t.name
}

func (t task) Next() string {
	return t.next
}

func (t task) Parameters() map[string]interface{} {
	return t.parameters
}
