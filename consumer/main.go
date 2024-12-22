package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var (
	rabbitMQURL = "amqp://stockmarket:supersecret123@rabbitmq:5672/" // URL zu RabbitMQ mit Benutzer
	queueName   = "AAPL"                                             // Der Name der Queue
)

// Struktur für die Nachricht
type StockMessage struct {
	Company   string  `json:"company"`
	EventType string  `json:"eventType"`
	Price     float64 `json:"price"`
}

func main() {
	// Verbindung zu RabbitMQ herstellen
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Verbindung zu RabbitMQ fehlgeschlagen: %v", err)
	}
	defer conn.Close()

	// Kanal öffnen
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Kanal konnte nicht geöffnet werden: %v", err)
	}
	defer ch.Close()

	// Queue deklarieren
	_, err = ch.QueueDeclare(
		queueName, // Name der Queue
		false,     // Durable: Queue überlebt Serverneustarts
		false,     // Auto-delete: Queue wird gelöscht, wenn sie nicht verwendet wird
		false,     // Exclusive: Die Queue ist nur für diese Verbindung
		false,     // No-wait: Kein Warten auf Bestätigung vom Server
		nil,       // Zusätzliche Argumente (optional)
	)
	if err != nil {
		log.Fatalf("Fehler beim Deklarieren der Queue: %v", err)
	}
	log.Println("Queue erfolgreich deklariert.")

	// Nachrichten aus der Queue konsumieren
	msgs, err := ch.Consume(
		queueName, // Queue-Name
		"",        // Consumer-Tag
		true,      // auto-acknowledge
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // zusätzliche Argumente
	)
	if err != nil {
		log.Fatalf("Fehler beim Konsumieren von Nachrichten: %v", err)
	}

	// Slice für gesammelte Nachrichten
	var data []float64
	batchSize := 1000

	// Nachrichten verarbeiten
	for msg := range msgs {
		// Nachricht in ein StockMessage-Objekt umwandeln
		var stockMsg StockMessage
		err := json.Unmarshal(msg.Body, &stockMsg)
		if err != nil {
			log.Printf("Fehler beim Parsen der JSON-Nachricht: %v", err)
			continue
		}

		// Den Preis aus der Nachricht extrahieren
		data = append(data, stockMsg.Price)

		// Wenn 1000 Nachrichten gesammelt sind, Durchschnitt berechnen
		if len(data) == batchSize {
			// Durchschnitt berechnen
			var sum float64
			for _, val := range data {
				sum += val
			}
			avg := sum / float64(batchSize)

			// Durchschnitt in der Konsole ausgeben
			fmt.Printf("Durchschnitt der letzten %d Nachrichten: %.2f\n", batchSize, avg)

			// Slice zurücksetzen für die nächste Gruppe
			data = nil
		}
	}
}
