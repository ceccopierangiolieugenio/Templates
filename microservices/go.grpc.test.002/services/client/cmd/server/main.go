package main

import (
	"log"
	"net/http"

	"test.90poe/services/client/pkg/api"
	"test.90poe/services/client/pkg/config"

)

var address string = "localhost:50051"

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	ln := r.URL.Path[1:]
	log.Printf("Handling: %s", ln)    

}

func main() {
	cfg, err := config.CreateConfig()
	if err != nil {
		log.Fatalf("failed initiating the config: %v", err)
	}

	pds, err := api.NewPDSClient(cfg.PDS_Host)
	if err != nil {
		log.Fatalf("failed initiating the PDS Client: %v", err)
	}
	defer pds.Close()

	go pds.ImportFile(cfg.InFile)

	api.RunWebServer(cfg.HTTP_Port, pds)
}