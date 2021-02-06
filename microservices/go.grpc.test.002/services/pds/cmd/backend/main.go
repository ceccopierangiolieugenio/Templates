package main

import (
	"os"
	"os/signal"
	"syscall"
	"log"
	"net"

	"google.golang.org/grpc"

	"test.90poe/services/pds/pkg/config"
	"test.90poe/services/pds/pkg/api"
	"test.90poe/services/proto_pds"
)

func main() {
	cfg, err := config.CreateConfig()
	if err != nil {
		log.Fatalf("failed initiating the config: %v", err)
	}

	/* Init gRPC */
	lis, err := net.Listen("tcp", cfg.PDS_Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("tcp listening at:", cfg.PDS_Port)
	
	gs := grpc.NewServer()
	srv, err := api.NewServer(cfg)
	if err != nil {
		log.Fatalf("failed to init the Server: %v", err)
	}

	proto_pds.RegisterPDSServer(gs, srv)

	/* Handle Syscalls */
	errChan := make(chan error)
	stopChan := make(chan os.Signal)

	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	
	go func() {
		log.Println("Staring gRPC...")
		if err := gs.Serve(lis); err != nil {
			errChan <- err
		}
	}()
	
	defer gs.GracefulStop()
	defer srv.CloseDbConnections()
		
	select {
		case err := <-errChan:
			log.Printf("Fatal error: %v\n", err) 
		case <-stopChan:
	}
}
