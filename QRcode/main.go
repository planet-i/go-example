package main

import (
	"log"

	"github.com/skip2/go-qrcode"
)

func main() {
	_, err := qrcode.Encode("https://example.org", qrcode.Medium, 256)
	if err != nil {
		log.Fatal(err)
	}

}
