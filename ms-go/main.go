package main

import (
	"ms-go/router"
	"ms-go/app/consumers"
)

func main() {
	go consumers.ConsumerKafka()

	router.Run()
}
