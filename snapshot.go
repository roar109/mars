package main

import (
	"log"
	"strings"
	"time"
)

func createConfigSnapshot() {
	t := time.Now()
	parsmStrs := []string{"config", t.Format("20060102150405"), "json"}
	newFileName := strings.Join(parsmStrs, ".")
	err := CopyFile("config.json", newFileName)

	if err != nil {
		log.Panic(err)
	}
}
