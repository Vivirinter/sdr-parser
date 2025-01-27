package main

import (
	"log"
	"github.com/Vivirinter/sdr-parser/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
