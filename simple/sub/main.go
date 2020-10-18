package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"os/signal"
	"time"
)

// $ go run ./simple/sub/main.go
// 2020/10/18 14:24:40 Connected
// 2020/10/18 14:25:16 Receive: Hello, World!
// 2020/10/18 14:25:19 Receive: Hello, World!
// 2020/10/18 14:25:22 Receive: Hello, World!
// 2020/10/18 14:25:25 Receive: Hello, World!
// 2020/10/18 14:25:28 Receive: Hello, World!
// 2020/10/18 14:25:31 Receive: Hello, World!
// ^C2020/10/18 14:25:34 Finish

func main() {
	// Setup client
	options := mqtt.NewClientOptions().
		AddBroker("tcp://localhost:1883").
		SetOnConnectHandler(func(_ mqtt.Client) {
			log.Println("Connected")
		}).
		SetConnectionLostHandler(func(_ mqtt.Client, err error) {
			log.Printf("Disconnected: %+v\n", err)
		}).
		SetClientID("Subscriber")
	c := mqtt.NewClient(options)

	// Establish connection
	token := c.Connect()
	if token.WaitTimeout(3 * time.Second) {
		err := token.Error()
		if err != nil {
			log.Fatalf("Failed to connect: %+v", err)
		}
	} else {
		log.Fatal("Timed out on connection establishment")
	}
	defer c.Disconnect(100)

	// Receive signal to finish operation
	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, os.Kill)

	// Subscribe to "test/topic" topic.
	token = c.Subscribe("test/topic", byte(0), func(_ mqtt.Client, message mqtt.Message) {
		log.Printf("Receive: %s\n", string(message.Payload()))
	})
	if token.WaitTimeout(3 * time.Second) {
		err := token.Error()
		if err != nil {
			log.Fatalf("Failed to subscribe: %+v", err)
		}
	} else {
		log.Fatalf("Timed out on subscription")
	}

	<-finish
	log.Println("Finish")
}
