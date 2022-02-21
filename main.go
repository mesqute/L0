package main

import (
	"L0/controllers"
	"L0/models"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"time"
)

func main() {

	clusterID := "test-cluster"
	clientID := "ThisIsMYClientId"
	// подключение к nats-streaming серверу
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("localhost:4222"))
	if err != nil {
		log.Fatal(err)
	}
	// подписка на канал и обработка сообщениий
	sub, err := sc.Subscribe("foo",
		func(m *stan.Msg) {
			// инициализация экземпляра структуры Order
			// в которую будут записаны полученные в запросе данные
			var order models.Order

			// декодирование полученного json файла в ранее созданый экземпяр структуры Order
			if err := json.Unmarshal(m.Data, &order); err != nil {
				log.Print(err)
				return
			}
			log.Print(order)
			// проверка полученных данных на корректность
			if err := order.Validate(); err != nil {
				log.Print(err)
				return
			}

			// добавление полученных данных в память сервиса
			if err := models.InsertData(order); err != nil {
				log.Print(err)
				return
			}
		},
		stan.DeliverAllAvailable(),
		stan.AckWait(20*time.Second),
		stan.MaxInflight(15))
	if err != nil {
		log.Fatal(err)
	}

	defer func(sub stan.Subscription) {
		err := sub.Unsubscribe()
		if err != nil {
			log.Fatal(err)
		}
	}(sub)

	server := &http.Server{
		Addr:    ":8080",
		Handler: controllers.GetRoutes(),
	}

	// запуск сервера
	log.Printf("Запуск сервера на %s", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Printf("[main] ListenAndServe: %s", err)
	}
}
