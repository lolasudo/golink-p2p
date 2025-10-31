package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	api "github.com/lolasudo/golink-p2p/practice-4/api"
)

func main() {
	// Загрузка TLS сертификатов для mTLS
	cert, err := tls.LoadX509KeyPair("../certs/client.pem", "../certs/client-key.pem")
	if err != nil {
		log.Fatalf("failed to load client cert: %v", err)
	}

	caCert, err := os.ReadFile("../certs/ca.pem")
	if err != nil {
		log.Fatalf("failed to read CA cert: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		ServerName:   "localhost",
	})

	// Подключаемся к серверу
	conn, err := grpc.Dial("127.0.0.1:9445", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := api.NewDiceServiceClient(conn)

	// Вызываем gRPC метод
	response, err := client.RollDie(context.Background(), &api.RollDieRequest{})
	if err != nil {
		log.Fatalf("failed to roll die: %v", err)
	}

	fmt.Printf("🎲 Rolled: %d\n", response.Value)
}
