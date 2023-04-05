package parser

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/AkronimBlack/process-manager/shared"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

const (
	StartNode = "start_node"
)

type Actions map[string]*Action

type Parser struct {
	handlers map[string]Handler
	actions  Actions
	sessions []Session

	lock sync.Mutex
}

func NewParser() *Parser {
	return &Parser{
		handlers: map[string]Handler{
			IsGreater:  IsGreaterHandler,
			IsLower:    IsLowerHandler,
			IsEqual:    IsEqualHandler,
			HttpAction: HttpHandler,
		},
		sessions: make([]Session, 0),
	}
}

func (p *Parser) runWebhook(session Session) {
	if session.OnFinishWebhook() == nil {
		session.OnFinishWebhookResponse()
		return
	}
	req, err := http.NewRequest(http.MethodPost, session.OnFinishWebhook().Url, bytes.NewBuffer(shared.ToJsonByte(session)))
	if err != nil {
		session.SetOnFinishWebhookResponse(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	client := &http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		session.SetOnFinishWebhookResponse(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		session.SetOnFinishWebhookResponse(map[string]interface{}{
			"status":      resp.Status,
			"status_code": resp.StatusCode,
			"error":       err.Error(),
		})
		return
	}
	session.SetOnFinishWebhookResponse(map[string]interface{}{
		"status":      resp.Status,
		"status_code": resp.StatusCode,
		"response":    string(body),
	})
}

func (p *Parser) SetActions(actions Actions) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.actions = actions
}

func (p *Parser) Actions() Actions {
	if p.actions == nil {
		return map[string]*Action{}
	}
	return p.actions
}

func (p *Parser) AddHandler(action string, handler Handler) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.handlers[action] = handler
}

func (p *Parser) ActionHandler(action string) Handler {
	v, ok := p.handlers[action]
	if !ok {
		return nil
	}
	return v
}

func (p *Parser) LoadFile(location string) error {
	if path.Ext(location) != ".json" {
		return fmt.Errorf("%s is not a json file", location)
	}
	file, err := os.ReadFile(location)
	if err != nil {
		return err
	}
	var actions Actions
	err = json.Unmarshal(file, &actions)
	if err != nil {
		return err
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	p.actions = actions
	return nil
}

func (p *Parser) Sessions() []Session {
	return p.sessions
}

type ValidateErrors map[string]ValidationErrors

func (e ValidateErrors) IsValid() bool {
	if len(e) == 0 {
		return true
	}
	return false
}

func (p *Parser) Validate() ValidateErrors {
	return p.validate(p.actions)
}

func (p *Parser) validate(actions Actions) ValidateErrors {
	errors := make(ValidateErrors)
	var hasStartNode bool
	for id, action := range actions {
		if action.ActionType == StartNode {
			hasStartNode = true
		}
		actionErrors := p.ValidateAction(action)
		if !actionErrors.IsValid() {
			errors[id] = actionErrors
		}
	}
	if !hasStartNode {
		errors[StartNode] = ValidationErrors{
			StartNode: {
				"start_node is required",
			},
		}
	}
	return errors
}

type ValidationErrors map[string][]string

func (e ValidationErrors) Add(key string, errors []string) {
	_, ok := e[key]
	if ok {
		e[key] = append(e[key], errors...)
		return
	}
	e[key] = errors
}

func (e ValidationErrors) Merge(errors ValidationErrors) {
	for key, value := range errors {
		e.Add(key, value)
	}
}

func (e ValidationErrors) IsValid() bool {
	if len(e) == 0 {
		return true
	}
	return false
}

func (p *Parser) ValidateAction(action *Action) ValidationErrors {
	errors := make(ValidationErrors, 0)
	if action.ActionType == StartNode {
		if action.OnSuccess == "" {
			errors.Add("on_success", []string{"on_success is a required field"})
		}
		return errors
	}

	if action.ActionType == "" {
		errors.Add("type", []string{"type is a required field"})
	}
	if action.OnSuccess == "" {
		errors.Add("on_success", []string{"on_success is a required field"})
	}
	if action.OnFailure == "" {
		errors.Add("on_failure", []string{"on_failure is a required field"})
	}
	return errors
}

func (p *Parser) Execute(ctx context.Context, data map[string]interface{}, webhook *Webhook) string {
	newSession := NewSession(data, webhook)
	p.sessions = append(p.Sessions(), newSession)
	startAction := p.actions[StartNode]
	firstAction := p.actions[startAction.OnSuccess]
	go p.runAction(ctx, firstAction, newSession)
	return newSession.Uuid()
}

func (p *Parser) runAction(ctx context.Context, action *Action, session Session) {
	log.Print(action)
	handler := p.ActionHandler(action.ActionType)
	if handler == nil {
		p.runWebhook(session)
		return
	}
	next := handler(ctx, action, session)
	if next == "" {
		p.runWebhook(session)
		return
	}
	p.runActionById(ctx, next, session)
}

func (p *Parser) runActionById(ctx context.Context, actionId string, session Session) {
	action := p.actions[actionId]
	if action == nil {
		p.runWebhook(session)
		return
	}
	p.runAction(ctx, action, session)
}
