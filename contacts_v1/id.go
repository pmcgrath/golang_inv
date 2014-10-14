package main

import (
	"crypto/rand"
	"fmt"
)

func Uuid() string {
	// See 	http://stackoverflow.com/questions/15130321/is-there-a-method-to-generate-a-uuid-with-go-language
	//	https://groups.google.com/forum/#!topic/golang-nuts/Rn13T6BZpgE
	//	https://groups.google.com/forum/#!msg/golang-nuts/d0nF_k4dSx4/rPGgfXv6QCoJ
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		return ""
	}
	buf[6] = (buf[6] & 0x0f) | 0x40
	buf[8] = (buf[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:])
}
