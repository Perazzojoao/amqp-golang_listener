package main

import "github.com/Perazzojoao/amqp-golang_listener/listener"

func main() {
	listenerConfig := listener.NewConfig()
	listener := listener.NewListener(listenerConfig)
	listener.Listen()
}
