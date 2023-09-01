package main

import (
	"fmt"
	"os"

	pgx "github.com/jackc/pgx/v5"
	stan "github.com/nats-io/stan.go"
)

type Subscriber struct {
	StanConn     stan.Conn
	Sub          stan.Subscription
	PostgresConn *pgx.Conn
}

func publisher(config NATSConfig, path string) error {
	sc, err := stan.Connect(config.Cluster, config.Publisher)
	if err != nil {
		return err
	}

	bytes, err := os.ReadFile("data/" + path)
	if err != nil {
		return err
	}
	err = sc.Publish(config.Channel, bytes)
	if err != nil {
		return err
	}
	sc.Close()
	return err
}

func subscriberStart(config NATSConfig, conn *pgx.Conn) (Subscriber, error) {
	var subscriber Subscriber
	var err error

	subscriber.StanConn, err = stan.Connect(config.Cluster, config.Subscriber)
	subscriber.PostgresConn = conn
	if err != nil {
		return subscriber, err
	}
	subscriber.Sub, err = subscriber.StanConn.Subscribe(config.Channel, func(message *stan.Msg) {
		// Если данные не являются валидным json-файлом,то они отбрасываются
		order, err := JsonToOrder(message.Data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid json: %v\n", err)
		} else {
			// При возникновении ошибки записи данные сохраняются в файл,
			// чтобы избежать потери данных
			err = InsertOrder(order, subscriber.PostgresConn)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Database write error: %v\n", err)
				os.WriteFile("backup.json", message.Data, os.ModeAppend)
				return
			}
			cache[order.OrderUID] = order
		}
	})
	return subscriber, err
}

func subscriberStop(subscriber Subscriber) error {
	err := subscriber.Sub.Unsubscribe()
	if err != nil {
		return err
	}

	err = subscriber.StanConn.Close()
	return err
}
