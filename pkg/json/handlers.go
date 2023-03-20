package json

import (
	"context"
	"fmt"
	"log"
	"strings"
)

const (
	IsGreater  = "is_greater"
	IsLower    = "is_lower"
	IsEqual    = "is_equal"
	HttpAction = "http"

	comparingKey = "comparing"
	compareToKey = "compare_to"

	result = "result"
)

func IsPlaceholder(value string) bool {
	if strings.Contains(value, "{{") && strings.Contains(value, "}}") {
		return true
	}
	return false
}

func CleanPlaceHolder(value string) string {
	return strings.TrimSpace(strings.TrimRight(strings.TrimLeft(value, "{{"), "}}"))
}

func IsGreaterHandler(ctx context.Context, action Action, session *Session) {
	resultVariable := action.Args.GetString(result, fmt.Sprintf("%s.result", action.ActionType))

	comparing := action.Args.GetString(comparingKey)
	comparingValue := session.PlaceholderOrIntValue(comparing)

	compareTo := action.Args.GetString(compareToKey)
	compareToValue := session.PlaceholderOrIntValue(compareTo)

	log.Println(comparingValue, compareToValue)

	session.Set(resultVariable, comparingValue > compareToValue)
}

func IsLowerHandler(ctx context.Context, action Action, session *Session) {
	resultVariable := action.Args.GetString(result, fmt.Sprintf("%s.result", action.ActionType))
	comparing := action.Args.Get(comparingKey)
	comparingValue := session.PlaceholderOrIntValue(comparing)

	compareTo := action.Args.Get(compareToKey)
	compareToValue := session.PlaceholderOrIntValue(compareTo)
	session.Set(resultVariable, comparingValue < compareToValue)
}

func IsEqualHandler(ctx context.Context, action Action, session *Session) {
	resultVariable := action.Args.GetString(result, fmt.Sprintf("%s.result", action.ActionType))

	comparing := action.Args.Get(comparingKey)
	comparingValue := session.PlaceholderOrIntValue(comparing)

	compareTo := action.Args.Get(compareToKey)
	compareToValue := session.PlaceholderOrIntValue(compareTo)

	session.Set(resultVariable, comparingValue == compareToValue)
}

func HttpHandler(ctx context.Context, action Action, session *Session) {

}
