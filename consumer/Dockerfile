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

# Falls erforderlich, portiere die Anwendung auf einen bestimmten Port (z.B. 8080)
EXPOSE 8080

# Führe die Go-Anwendung aus
#CMD ["./consumer"]
