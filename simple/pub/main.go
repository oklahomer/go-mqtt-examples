package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"os/signal"
	"time"
)

// $ go run ./simple/pub/main.go
// 2020/10/18 14:25:13 Connected
// 2020/10/18 14:25:16 Published
// 2020/10/18 14:25:19 Published
// 2020/10/18 14:25:22 Published
// 2020/10/18 14:25:25 Published
// 2020/10/18 14:25:28 Published
// 2020/10/18 14:25:31 Published
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
		SetClientID("Publisher")
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

	// Operate every 3 seconds
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// Run til signal is given
	for {
		select {
		case <-ticker.C:
			// Publish
			token := c.Publish("test/topic", byte(0), false, []byte("Hello, World!"))
			if token.WaitTimeout(1 * time.Second) {
				err := token.Error()
				if err != nil {
					log.Printf("Failed to publish: %+v\n", err)
					continue
				}

				log.Println("Published")
				continue
			}

			log.Println("Timed out")

		case <-finish:
			log.Println("Finish")
			return
		}
	}
}
