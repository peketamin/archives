// ref: [Go Playground](http://play.golang.org/p/OamEBLUljP)
package main

import (
	"fmt"
	"regexp"
)

func main() {
	src := []byte("eo eo eo eo")
	search := regexp.MustCompile("e")
	repl := []byte("AEI")

	i := 0
	src = search.ReplaceAllFunc(src, func(s []byte) []byte {
		if i < 2 {
			i += 1
			return repl
		}
		return s
	})

	fmt.Println(string(src))
}
