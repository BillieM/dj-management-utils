package helpers

func ContainsNonEmptyString(s []string) bool {
	for _, v := range s {
		if v != "" {
			return true
		}
	}
	return false
}
