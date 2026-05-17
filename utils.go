package main

import (
	"log"
	"time"
)

func MeasureExecutionTime(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took %s.", name, time.Since(start))
	}
}
