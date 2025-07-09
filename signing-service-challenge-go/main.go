package main

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
)

const (
	ListenAddress = ":8080"
	// TODO: add further configuration parameters here ...
)

func main() {
	db := persistence.NewInMemoryDatabase()
	deviceService := domain.NewDeviceService(db)
	server := api.NewServer(ListenAddress, deviceService)

	log.Printf("Listening on %s\n", ListenAddress)
	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
