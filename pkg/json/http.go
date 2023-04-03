package json

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

func (p *ParserHttpHandler) StartSession(ctx *gin.Context) {
	p.parser.Execute(context.Background())
}
