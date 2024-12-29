package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	rabbitMQURL = "amqp://stockmarket:supersecret123@rabbitmq:5672/"                                                           // URL zu RabbitMQ mit Benutzer
	queueName   = "AAPL"                                                                                                       // Der Name der Queue
	mongoURI    = "mongodb://host.docker.internal:27017,host.docker.internal:27018,host.docker.internal:27019/?replicaSet=rs0" // URI zu MongoDB
)

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

	// Verbindung zur MongoDB aufbauen
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Fehler beim Verbinden mit MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// MongoDB-Datenbank auswählen oder erstellen
	database := client.Database("stockmarket")

	// Collection "stocks" auswählen
	collection := database.Collection("stocks")

	// Sicherstellen, dass die Collection mit einem Index erstellt wird
	// Hier können wir optional einen Index hinzufügen, falls erforderlich
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"company": 1, // Beispiel: Index auf das 'company'-Feld
		},
		Options: nil,
	}

	// Erstellen des Indexes (falls notwendig)
	_, err = collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Printf("Fehler beim Erstellen des Indexes: %v", err)
	}

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

	// Nachrichten in Batches von 1000 abholen
	var batch []StockMessage
	batchSize := 1000
	for msg := range msgs {
		// Nachricht in ein StockMessage-Objekt umwandeln
		var stockMsg StockMessage
		err := json.Unmarshal(msg.Body, &stockMsg)
		if err != nil {
			log.Printf("Fehler beim Parsen der JSON-Nachricht: %v", err)
			continue
		}

		// Nachricht zur Batch hinzufügen
		batch = append(batch, stockMsg)

		// Wenn Batch die gewünschte Größe erreicht hat
		if len(batch) >= batchSize {
			// Durchschnitt berechnen
			var sum float64
			for _, msg := range batch {
				sum += msg.Price
			}
			avgPrice := sum / float64(len(batch))

			// Dokument für MongoDB erstellen
			document := bson.D{
				{Key: "company", Value: batch[0].Company}, // Alle Nachrichten im Batch haben denselben "company"-Wert
				{Key: "avgPrice", Value: avgPrice},
				{Key: "timestamp", Value: time.Now()},
			}

			// Dokument in die MongoDB-Sammlung einfügen
			_, err := collection.InsertOne(context.Background(), document)
			if err != nil {
				log.Printf("Fehler beim Einfügen der Nachricht in MongoDB: %v", err)
			} else {
				log.Printf("Durchschnittspreis für %s berechnet und in MongoDB eingefügt: %.2f", batch[0].Company, avgPrice)
			}

			// Batch zurücksetzen
			batch = nil
		}
	}
}
