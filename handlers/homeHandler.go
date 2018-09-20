package handlers

import (
	"net/http"
	"log"
	"time"
	"fmt"

	"github.com/OtgonbayarT/microservice/models"
)

type HandlersLog struct {
	logger *log.Logger
	dbUrl string
}


func (h *HandlersLog) EncodeHandler(w http.ResponseWriter, r *http.Request){
	url := r.PostFormValue("url")

	shortUrl, err := models.InsertUrl(h.dbUrl, url)
	if err != nil {
		h.logger.Printf("cannot save data: %v", err)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Cannot save url at the moment!"))
		return
	}

	w.WriteHeader(http.StatusOK)	
	w.Write([]byte(fmt.Sprintf("%s/%s/%s", r.Host, "decode", shortUrl)))
}

func (h *HandlersLog) DecodeHandler(w http.ResponseWriter, r *http.Request){
	shortUrl := r.URL.Path[len("/decode/"):]

	longUrl, err := models.GetUrl(h.dbUrl, shortUrl)
	if err != nil {
		h.logger.Printf("cannot save data: %v", err)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("we cannot find encoded url"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(longUrl))

}

func (h *HandlersLog) RedirectHandler(w http.ResponseWriter, r *http.Request){
	shortUrl := r.URL.Path[len("/redirect/"):]

	longUrl, err := models.GetUrl(h.dbUrl, shortUrl)
	if err != nil {
		h.logger.Printf("cannot save data: %v", err)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("we cannot find encoded url for redirect"))
		return
	}

	http.Redirect(w, r, string(longUrl), 301)
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
