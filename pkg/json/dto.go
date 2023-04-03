package json

func NewSessionDto(session *Session) SessionDto {
	if session == nil {
		return SessionDto{}
	}
	return SessionDto{
		Values:          session.values,
		ExecutedActions: session.executedActions,
	}
}

type SessionDto struct {
	Values          map[string]interface{} `json:"values"`
	ExecutedActions []*Action              `json:"executed_actions"`
}
