package json

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
