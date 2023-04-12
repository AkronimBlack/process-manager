package parser

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MessageResponse struct {
	Message string `json:"message"`
}

type ParserHttpHandler struct {
	parser *Parser
}

func NewParserHttpHandler(file string) ParserHttpHandler {
	parser := NewParser()
	if file != "" {
		err := parser.LoadFile(file)
		if err != nil {
			panic(err)
		}
	}
	return ParserHttpHandler{
		parser: parser,
	}
}

func BuildHttp(router *gin.Engine, file string) {
	httpHandler := NewParserHttpHandler(file)
	router.Group("api").
		GET("/sessions", httpHandler.GetSessions).
		GET("/sessions/:id", httpHandler.Session).
		POST("/sessions", httpHandler.StartSession).
		GET("/sessions/:id/tasks", httpHandler.Tasks).
		POST("/sessions/:id/tasks/:task_id", httpHandler.CompleteTask)
}

func (p *ParserHttpHandler) GetSessions(ctx *gin.Context) {
	sessions := p.parser.Sessions()
	dtoSessions := make([]SessionDto, len(sessions))
	for i, activeSession := range sessions {
		dtoSessions[i] = NewSessionDto(activeSession)
	}
	ctx.JSON(http.StatusOK, dtoSessions)
}

type StartSessionRequest struct {
	Data    map[string]interface{}     `json:"data"`
	Webhook StartSessionWebhookRequest `json:"webhook"`
}

type StartSessionWebhookRequest struct {
	Url string `json:"url"`
}

func (p *ParserHttpHandler) StartSession(ctx *gin.Context) {
	var request StartSessionRequest
	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(
			http.StatusUnprocessableEntity,
			MessageResponse{Message: err.Error()},
		)
		return
	}

	sessionUuid := p.parser.Execute(context.Background(), request.Data, NewWebHook(request.Webhook.Url))
	ctx.JSON(
		http.StatusCreated,
		MessageResponse{Message: sessionUuid},
	)
}

func (p *ParserHttpHandler) Session(ctx *gin.Context) {
	activeSession := p.parser.Session(ctx.Param("id"))
	if activeSession == nil {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}
	ctx.JSON(http.StatusOK, NewSessionDto(activeSession))
}

func (p *ParserHttpHandler) Tasks(ctx *gin.Context) {
	activeSession := p.parser.Session(ctx.Param("id"))
	if activeSession == nil {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}
	ctx.JSON(http.StatusOK, NewTasksDto(activeSession.Tasks()))
}

func (p *ParserHttpHandler) CompleteTask(ctx *gin.Context) {
	activeSession := p.parser.Session(ctx.Param("id"))
	if activeSession == nil {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}
	activeTask := activeSession.Task(ctx.Param("task_id"))
	if activeTask == nil {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}
	var payload map[string]interface{}
	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	activeTask.Execute(payload)
	ctx.JSON(http.StatusOK, NewTaskDto(activeTask))
}
