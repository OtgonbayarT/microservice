package handlers

import (
	"net/http"
	"log"
	"time"
	"fmt"

	"github.com/OtgonbayarT/microservice/controllers"
	"github.com/rapidloop/skv"
)

const message = "home handler!!"

type HandlersLog struct {
	logger *log.Logger
	dbUrl string
}


func (h *HandlersLog) EncodeHandler(w http.ResponseWriter, r *http.Request){
	url := r.PostFormValue("url")

	store, err := skv.Open(h.dbUrl)
	if err!= nil {
		store.Close()
		h.logger.Printf("cannot open db: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}

	if err := store.Put(fmt.Sprint(controllers.Hash(url)), url); err != nil {
		store.Close()
		h.logger.Printf("cannot save data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}

	store.Close()

	w.WriteHeader(http.StatusOK)	
	w.Write([]byte(fmt.Sprint(controllers.Hash(url))))
}

func (h *HandlersLog) DecodeHandler(w http.ResponseWriter, r *http.Request){
	code := r.URL.Path[len("/decode/"):]

	store, dberr := skv.Open(h.dbUrl)
	if dberr!= nil {
		store.Close()
		h.logger.Printf("cannot open db: %v", dberr)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}

	var val string
	if err := store.Get(fmt.Sprint(code), &val); err != nil {
		store.Close()
		h.logger.Printf("data not found: %v", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("URL Not Found."))
		return
	}

	store.Close()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(val))

}

func (h *HandlersLog) RedirectHandler(w http.ResponseWriter, r *http.Request){
	code := r.URL.Path[len("/redirect/"):]

	store, dberr := skv.Open(h.dbUrl)
	if dberr!= nil {
		store.Close()
		h.logger.Printf("cannot open db: %v", dberr)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}

	var url string
	if err := store.Get(fmt.Sprint(code), &url); err != nil {
		store.Close()
		h.logger.Printf("data not found: %v", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("URL Not Found."))
		return
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


func NewHandlersLog(logger *log.Logger, dbUrl string ) *HandlersLog{
	return &HandlersLog{
		logger: logger,
		dbUrl: dbUrl,
	}
}
