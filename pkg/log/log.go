package log

import (
	"fmt"
	"os"
)

func Error(v any) {
	//nolint:forbidigo //dn
	fmt.Printf("%v\n", v)
	os.Exit(1)
}

func Info(v any) {
	//nolint:forbidigo //dn
	fmt.Printf("%v\n", v)
}
