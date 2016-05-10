package main

import (
	"log"
	"github.com/benmanns/goworker"
	"./task"
)

func main () {
	if err := goworker.Work(); err != nil {
		log.Println("Error:", err)
	}
}

func init () {
	log.Println("Start worker")
	goworker.Register("GetPageItems", task.GetPageItems)
	goworker.Register("GetItem", task.GetItem)
}