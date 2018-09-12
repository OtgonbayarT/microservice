package handlers

import (
	"net/http"
	"log"
	"time"
)

const message = "home handler!!"

type HandlersLog struct {
	logger *log.Logger
}

func (h *HandlersLog) HomeHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func (h *HandlersLog) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		startTime := time.Now()
		defer h.logger.Printf("request processed in %s\n", time.Now().Sub(startTime));
		next(w, r)
	}
}


func (h *HandlersLog) SetUpRoutes(mux *http.ServeMux){
	mux.HandleFunc("/", h.Logger(h.HomeHandler))
}

func NewHandlersLog(logger *log.Logger) *HandlersLog{
	return &HandlersLog{
		logger: logger,
	}
}
