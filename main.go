package main

import (
	"net/http"
	"log"
	"os"

	"github.com/OtgonbayarT/microservice/server"
	"github.com/OtgonbayarT/microservice/handlers"
	"github.com/prometheus/client_golang/prometheus"
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

	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "encode_duration_seconds",
		Help: "Time taken to create url",
	}, []string{"code"})	
	
	h := handlers.NewHandlersLog(logger, dbUrl, histogram)

	prometheus.Register(histogram)
	
	mux := http.NewServeMux()	
	h.SetUpRoutes(mux)
	srv := server.New(mux, MSServiceAddr)
	
	logger.Println("server starting")
	err := srv.ListenAndServe()
	if err!= nil {
		logger.Fatalf("server failed to start: %v", err)
	}
	
}
