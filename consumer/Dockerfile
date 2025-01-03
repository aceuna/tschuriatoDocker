# Wähle ein offizielles Golang-Image als Basis
FROM golang:1.20-alpine

# Setze das Arbeitsverzeichnis im Container
WORKDIR /app

# Kopiere die Go Modul-Dateien und lade die Abhängigkeiten
COPY go.mod go.sum ./
RUN go mod tidy

# Kopiere den restlichen Code
COPY . .

# Baue die Go-Anwendung
RUN go build -o consumer .