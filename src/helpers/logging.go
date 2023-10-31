package helpers

import "fmt"

func HandleFatalError(err error) {
	panic(err)
}

func WriteToLog(message string) {
	fmt.Println(message)
}
