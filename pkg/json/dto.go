package json

func NewSessionDto(session *Session) SessionDto {
	if session == nil {
		return SessionDto{}
	}
	return SessionDto{
		Uuid:            session.Uuid,
		Values:          session.values,
		ExecutedActions: session.executedActions,
	}
}

type SessionDto struct {
	Uuid            string                 `json:"uuid"`
	Values          map[string]interface{} `json:"values"`
	ExecutedActions []*Action              `json:"executed_actions"`
}
