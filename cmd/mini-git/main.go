package main

import (
	"crypto/sha1"
	"fmt"
	"io"
)

func main() {
	h := sha1.New()
	io.WriteString(h, "hello world")
	fmt.Printf("%x\n", h.Sum(nil))
}
