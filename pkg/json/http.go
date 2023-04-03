package json

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ParserHttpHandler struct {
	parser *Parser
}

func AttachHttpHandlers(router *gin.Engine) {
	httpHandler := ParserHttpHandler{
		parser: NewParser(),
	}
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

func (p *ParserHttpHandler) StartSession(ctx *gin.Context) {
	p.parser.Execute(context.Background())
}
