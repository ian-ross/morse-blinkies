package main

import (
	"os"
	"strings"

	"github.com/ian-ross/morse-blinkies/mbaas/processing"
)

func main() {
	s := strings.Join(os.Args[1:], " ")
	processing.MakeBoard(s, 0)
}
