package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"gitlab.com/back1ng1/messaggio-test/internal/entity"
)

var topic = "messages"
var partition = 0

type Kafka struct {
	conn *kafka.Conn
}

func NewConn() Kafka {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "kafka:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	return Kafka{
		conn: conn,
	}
}

func NewReader() <-chan kafka.Message {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"kafka:9092"},
		Topic:     topic,
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})

	msgs := make(chan kafka.Message)

	go func() {
		for {
			msg, err := r.ReadMessage(context.Background())
			if err != nil {
				close(msgs)

				if err := r.Close(); err != nil {
					log.Fatal("failed to close reader:", err)
				}

				break
			}

			msgs <- msg
		}
	}()

	return msgs
}

func (k Kafka) StoreMessage(msg entity.Message) error {
	json, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Value: json,
	}

	_, err = k.conn.WriteMessages(message)
	if err != nil {
		return err
	}

	return nil
}
