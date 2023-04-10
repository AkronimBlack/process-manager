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
		POST("/sessions", httpHandler.StartSession)
}

func (p *ParserHttpHandler) GetSessions(ctx *gin.Context) {
	sessions := p.parser.Sessions()
	dtoSessions := make([]SessionDto, len(sessions))
	for i, session := range sessions {
		dtoSessions[i] = NewSessionDto(session)
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
