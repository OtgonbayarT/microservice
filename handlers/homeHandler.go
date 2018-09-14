package handlers

import (
	"net/http"
	"log"
	"time"
	"fmt"
	"errors"

	"github.com/OtgonbayarT/microservice/controllers"
	"github.com/rapidloop/skv"
)

const message = "home handler!!"

type HandlersLog struct {
	logger *log.Logger
}

var (
	ErrNotFound = errors.New("skv: key not found")
)

func (h *HandlersLog) EncodeHandler(w http.ResponseWriter, r *http.Request){
	url := r.PostFormValue("url")
	store, err := skv.Open("/home/otgonbayar/urls.db")
	if err != nil {
		h.logger.Fatalf("cannot open db: %v", err)
	}
	if err := store.Put(fmt.Sprint(controllers.Hash(url)), url); err != nil {
		h.logger.Fatalf("cannot save data: %v", err)
	}
	store.Close()
	w.WriteHeader(http.StatusOK)
	
	w.Write([]byte(fmt.Sprint(controllers.Hash(url))))
}

func (h *HandlersLog) DecodeHandler(w http.ResponseWriter, r *http.Request){
	code := r.URL.Path[len("/decode/"):]
	store, err := skv.Open("/home/otgonbayar/urls.db")
	if err != nil {
		h.logger.Fatalf("cannot open db: %v", err)
	}
	var val string
	if err := store.Get(fmt.Sprint(code), &val); err != nil {
		h.logger.Fatalf("data not found: %v", err)
	}
	store.Close()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(val))
}

func (h *HandlersLog) RedirectHandler(w http.ResponseWriter, r *http.Request){
	code := r.URL.Path[len("/redirect/"):]
	store, err := skv.Open("/home/otgonbayar/urls.db")
	if err != nil {
		h.logger.Fatalf("cannot open db: %v", err)
	}
	var url string
	if err := store.Get(fmt.Sprint(code), &url); err != nil {
		h.logger.Fatalf("data not found: %v", err)
	}
	store.Close()

	http.Redirect(w, r, string(url), 301)
}

func (h *HandlersLog) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		startTime := time.Now()
		defer h.logger.Printf("request processed in %s\n", time.Now().Sub(startTime));
		next(w, r)
	}
}


func (h *HandlersLog) SetUpRoutes(mux *http.ServeMux){
	mux.HandleFunc("/encode", h.Logger(h.EncodeHandler))
	mux.HandleFunc("/decode/", h.Logger(h.DecodeHandler))
	mux.HandleFunc("/redirect/", h.Logger(h.RedirectHandler))
}


func NewHandlersLog(logger *log.Logger) *HandlersLog{
	return &HandlersLog{
		logger: logger,
	}
}
