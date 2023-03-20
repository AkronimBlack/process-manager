package json

const (
	Url     = "url"
	Method  = "method"
	Timeout = "timeout"
	Headers = "headers"
	Payload = "payload"
)

// ValidateTwoValueOperators validates is_greater,is_lower,is_equal
func ValidateTwoValueOperators(action Action) ValidationErrors {
	errors := make(ValidationErrors, 0)
	comparing := action.Args.GetString(comparingKey)
	if comparing == "" {
		errors.Add(comparingKey, []string{"is mandatory"})
	}
	compareTo := action.Args.GetString(compareToKey)
	if compareTo == "" {
		errors.Add(compareToKey, []string{"is mandatory"})
	}
	return errors
}

func ValidateHttpAction(action Action) ValidationErrors {
	errors := make(ValidationErrors, 0)
	stringArgs := []string{
		Url,
		Method,
		result,
	}

	intArgs := []string{
		Timeout,
	}

	mapArgs := []string{
		Headers,
		Payload,
	}

	for _, arg := range stringArgs {
		value := action.Args.GetString(arg)
		if value == "" {
			errors.Add(arg, []string{"is mandatory"})
		}
	}

	for _, arg := range mapArgs {
		value := action.Args.GetMap(arg)
		if value == nil {
			errors.Add(arg, []string{"is mandatory"})
		}
	}

	for _, arg := range intArgs {
		value := action.Args.GetInt(arg)
		if value == 0 {
			errors.Add(arg, []string{"is mandatory"})
		}
	}

	return errors
}
