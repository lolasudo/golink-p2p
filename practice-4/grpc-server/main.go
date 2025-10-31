package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	api "github.com/lolasudo/golink-p2p/practice-4/api"
)

type diceServer struct {
	api.UnimplementedDiceServiceServer
	rand *rand.Rand
}

func (s *diceServer) RollDie(ctx context.Context, req *api.RollDieRequest) (*api.RollDieResponse, error) {
	// Генерируем случайное число от 1 до 6
	value := s.rand.Intn(6) + 1
	log.Printf("Rolled die: %d", value)
	return &api.RollDieResponse{Value: int32(value)}, nil
}

func main() {
	// Инициализируем генератор случайных чисел
	source := rand.NewSource(time.Now().UnixNano())
	server := &diceServer{
		rand: rand.New(source),
	}

	// Загрузка TLS сертификатов
	cert, err := tls.LoadX509KeyPair("../certs/server.pem", "../certs/server-key.pem")
	if err != nil {
		log.Fatalf("failed to load server cert: %v", err)
	}

	caCert, err := os.ReadFile("../certs/ca.pem")
	if err != nil {
		log.Fatalf("failed to read CA cert: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})

	// Создаем gRPC сервер с TLS
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	api.RegisterDiceServiceServer(grpcServer, server)

	// Запускаем сервер
	lis, err := net.Listen("tcp", "127.0.0.1:8443")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("gRPC server running on 127.0.0.1:8443")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
