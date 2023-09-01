package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	pgx "github.com/jackc/pgx/v5"
)

// Данные для подключения к NATS streaming
type NATSConfig struct {
	Cluster    string `json:"cluster"`
	Subscriber string `json:"subscriber"`
	Publisher  string `json:"publisher"`
	Channel    string `json:"channel"`
}

// Данные для подключения к Postgres
type PostgresConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Database string `json:"database"`
}

// Структура, в которую считываются данные конфигурации из data/config.json
type Config struct {
	Postgres PostgresConfig `json:"postgres"`
	NATS     NATSConfig     `json:"NATS"`
}

func main() {
	var config Config
	var wg sync.WaitGroup

	// Чтение конфигурации
	configBytes, err := os.ReadFile("data/config.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read config: %v\n", err)
		os.Exit(1)
	}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		os.Exit(1)
	}

	// Если программа запущена с двумя аргументами, то
	// Первый аргумент всегда `--publisher`, второй - путь к json с данными
	if len(os.Args[1:]) == 2 {
		if os.Args[1] == "--publisher" {
			publisher(config.NATS, os.Args[2])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to connect to NATS: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Unknown args: %s\n", os.Args[1])
		}
		return
	}

	// Данные для подключения к БД читаются из структуры
	postgresUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.Postgres.Username, config.Postgres.Password,
		config.Postgres.Address, config.Postgres.Port,
		config.Postgres.Database)

	conn, err := pgx.Connect(context.Background(), postgresUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Чтение таблиц из базы в кэш
	err = LoadDatabase(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load to database: %v\n", err)
		os.Exit(1)
	}

	// Подключение к NATS с передачей конфигурации
	var subscriber Subscriber
	subscriber, err = subscriberStart(config.NATS, conn)
	defer subscriberStop(subscriber)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to NATS: %v\n", err)
		os.Exit(1)
	}

	// Подключение веб-сервера
	wg.Add(1)
	go func() {
		StartHttp()
		defer wg.Done()
	}()

	wg.Wait()
}
