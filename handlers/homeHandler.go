package handlers

import (
	"net/http"
	"log"
	"time"
	"fmt"

	"github.com/OtgonbayarT/microservice/controllers"
)

const message = "home handler!!"

type HandlersLog struct {
	logger *log.Logger
}

func (h *HandlersLog) HomeHandler(w http.ResponseWriter, r *http.Request){
	// w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	url := r.PostFormValue("url")
	w.WriteHeader(http.StatusOK)
	// w.Write([]byte(url))
	w.Write([]byte(fmt.Sprint(controllers.Hash(url))))
}

func (h *HandlersLog) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		startTime := time.Now()
		defer h.logger.Printf("request processed in %s\n", time.Now().Sub(startTime));
		next(w, r)
	}
}


func (h *HandlersLog) SetUpRoutes(mux *http.ServeMux){
	mux.HandleFunc("/encode", h.Logger(h.HomeHandler))
	mux.HandleFunc("/decode", h.Logger(h.HomeHandler))
	mux.HandleFunc("/redirect", h.Logger(h.HomeHandler))
}


func NewHandlersLog(logger *log.Logger) *HandlersLog{
	return &HandlersLog{
		logger: logger,
	}
}
