package main

// CheckRole checks if the participant's role allows the requested action
func CheckRole(role string, action string) bool {
	rolePermissions := map[string][]string{
		"first responder":     {"create", "delete", "display", "transfer"},
		"second investigator": {"create", "delete", "display", "transfer"},
		"prosecutor":          {"display", "transfer"},
		"defense":             {"display", "transfer"},
		"court":               {"display", "transfer"},
	}

	allowedActions, exists := rolePermissions[role]
	if !exists {
		return false
	}

	for _, allowed := range allowedActions {
		if allowed == action {
			return true
		}
	}
	return false
}
