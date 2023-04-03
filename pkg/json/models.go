package json

import (
	"context"
	"encoding/json"
	"github.com/AkronimBlack/process-manager/shared"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"strconv"
	"sync"
)

// Handler interface. Expects id of next action to execute. Returning empty string finished the process execution
type Handler func(ctx context.Context, action *Action, session *Session) string

type Args map[string]interface{}

func (a Args) Get(key string) interface{} {
	return a[key]
}

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

type Session struct {
	Uuid            string
	values          map[string]interface{}
	executedActions []*Action
	inputData       map[string]interface{}
	lock            sync.Mutex
}

func NewSession(data map[string]interface{}) *Session {
	return &Session{
		Uuid:            uuid.NewString(),
		values:          make(map[string]interface{}),
		executedActions: make([]*Action, 0),
		inputData:       data,
	}
}

func (s *Session) ValueOf(key string) interface{} {
	value := gjson.Get(shared.ToJsonString(s.values), key)
	if value.Exists() {
		return value.Value()
	}
	return nil
}

func (s *Session) StringValueOf(key string, defaultValue string) string {
	value := gjson.Get(shared.ToJsonString(s.values), key)
	if value.Exists() {
		return value.String()
	}
	return defaultValue
}

func (s *Session) IntValueOf(key string, defaultValue int64) int64 {
	value := gjson.Get(shared.ToJsonString(s.values), key)
	if value.Exists() {
		return value.Int()
	}
	return defaultValue
}

func (s *Session) PlaceholderOrStringValue(value string) string {
	if IsPlaceholder(value) {
		return s.StringValueOf(CleanPlaceHolder(value), value)
	}
	return value
}

func (s *Session) PlaceholderOrIntValue(value interface{}) int64 {
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

func (s *Session) Set(key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.values[key] = value
}
