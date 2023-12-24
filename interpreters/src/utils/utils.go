package utils

import (
	"fmt"
	"os"
)

func Check(err error) {
	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(-1)
	}
}
