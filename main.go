package main

import (
	"net/http"
	"log"
	"os"

	"github.com/OtgonbayarT/microservice/server"
	"github.com/OtgonbayarT/microservice/handlers"
)

var (
	//MSCertFile     =  os.Getenv("MS_CERT_FILE")
	//MSKeyFile      =  os.Getenv("MS_KEY_FILE")
	//MSServiceAddr  =  os.Getenv("MS_SERVICE_ADDR")
	MSServiceAddr  =  ":8080"
)

func main() {
	logger := log.New(os.Stdout, "url shortener service: ", log.LstdFlags | log.Lshortfile)

	dbUrl  := "./urls.db"
	
	h := handlers.NewHandlersLog(logger, dbUrl)
	
	mux := http.NewServeMux()	
	h.SetUpRoutes(mux)
	srv := server.New(mux, MSServiceAddr)
	
	logger.Println("server starting")
	err := srv.ListenAndServe()
	if err!= nil {
		logger.Fatalf("server failed to start: %v", err)
	}
	
}
