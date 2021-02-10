package utils

import (
	"log"
	"runtime/debug"
)

func CheckError(e error) {
	if e != nil {
		log.Fatal(e, string(debug.Stack()))
	}
}
