package util

import (
	"fmt"
	"os"
)

const START_FAILED = 1

func CheckError(text string, err error) {
	if err == nil {return}
	_, _ = fmt.Fprintf(os.Stderr, "error: %s, err: %s\n", text, err)
	os.Exit(START_FAILED)
}
