package main

import (
	"github.com/nats-io/stan.go"
	"io/ioutil"
	"log"
)

func main() {
	clusterID := "test-cluster"
	clientID := "ThisIsMY"
	// подключение к nats-streaming серверу
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("localhost:4222"))
	if err != nil {
		log.Fatal(err)
	}
	dir, err := ioutil.ReadDir("./natsclientsender/json/")
	if err != nil {
		log.Print(err)
		return
	}
	for _, info := range dir {
		bytes, err := ioutil.ReadFile("./natsclientsender/json/" + info.Name())
		if err != nil {
			log.Print(err)
			return
		}
		err = sc.Publish("foo", bytes)
		if err != nil {
			log.Print(err)
			return
		}
	}
}
