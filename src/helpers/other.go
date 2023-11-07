package helpers

import (
	"encoding/json"
	"fmt"
)

func ContainsNonEmptyString(s []string) bool {
	for _, v := range s {
		if v != "" {
			return true
		}
	}
	return false
}

func PrettyPrint(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "\t")
	fmt.Println(string(s))
}
