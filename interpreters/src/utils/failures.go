package utils

import (
	"fmt"
	"os"
	"strconv"
)

func ReportFailure(err error) {
	loxError, isType := err.(LoxError)

	if isType {
		fmt.Fprintln(os.Stderr, "[line "+strconv.Itoa(loxError.Line)+"] Error "+loxError.Where+": "+loxError.Message)
	} else {
		fmt.Fprintln(os.Stderr, "Unexpected error: "+err.Error())
	}
}
