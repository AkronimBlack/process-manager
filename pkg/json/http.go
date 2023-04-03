package json

import (
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
		GET("/sessions", httpHandler.GetSessions)
}

func (p *ParserHttpHandler) GetSessions(ctx *gin.Context) {
	sessions := p.parser.Sessions()
	dtoSessions := make([]SessionDto, len(sessions))
	for i, session := range sessions {
		dtoSessions[i] = NewSessionDto(session)
	}
	ctx.JSON(http.StatusOK, dtoSessions)
}
