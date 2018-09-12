package main

import (
	"net/http"
	"log"
	"os"

	"github.com/OtgonbayarT/microservice/server"
	"github.com/OtgonbayarT/microservice/handlers"
)

const message = "hello world http service"

var (
	//GcukCertFile     =  os.Getenv("GCUK_CERT_FILE")
	//GcukKeyFile      =  os.Getenv("GCUK_KEY_FILE")
	//GcukServiceAddr  =  os.Getenv("GCUK_SERVICE_ADDR")
	GcukServiceAddr  =  ":8080"
)

func main() {
    logger := log.New(os.Stdout, "gcuk ", log.LstdFlags | log.Lshortfile)
	
	h := handlers.NewHandlersLog(logger)
	
	mux := http.NewServeMux()	
	h.SetUpRoutes(mux)
	srv := server.New(mux, GcukServiceAddr)
	
	logger.Println("server starting")
	err := srv.ListenAndServe()
	if err!= nil {
		logger.Fatalf("server failed to start: %v", err)
	}
	
}
