package main

import (
	"fmt"
	"strings"
)

func main() {

	r := map[string]string{
		"a": "1",
	}

	fmt.Println(r["p"])

}

func SubStr(value string, separator string) string {
	p := strings.Index(value, separator)
	if p > -1 {
		return value[0:p]
	}
	return value
}

func SubStrRight(value string, separator string) string {
	p := strings.LastIndex(value, separator)
	if p > -1 {
		return value[0:p]
	}
	return value
}
