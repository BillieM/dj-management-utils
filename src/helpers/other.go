package helpers

func ContainsNonEmptyString(s []string) bool {
	for _, v := range s {
		if v != "" {
			return true
		}
	}
	return false
}

func CountNewLines(s string) int {
	return countRune(s, '\n')
}

func countRune(s string, r rune) int {
	count := 0
	for _, c := range s {
		if c == r {
			count++
		}
	}
	return count
}
