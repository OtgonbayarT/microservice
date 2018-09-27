package handlers

import (
	"net/http"
	"log"
	"time"
	"fmt"

	"github.com/OtgonbayarT/microservice/models"
	"github.com/prometheus/client_golang/prometheus"
)

type HandlersLog struct {
	logger *log.Logger
	dbUrl string
	histogram *prometheus.HistogramVec
}


func (h *HandlersLog) EncodeHandler(w http.ResponseWriter, r *http.Request){
	contentType := r.Header.Get("Content-type")
    
	url := r.PostFormValue("url")

	if (contentType == "application/x-www-form-urlencoded" && len(url) > 0){
		shortUrl, err := models.InsertUrl(h.dbUrl, url)
		if err != nil {
			h.logger.Printf("cannot save data: %v", err)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Cannot save url at the moment!"))
			return
		}

		w.WriteHeader(http.StatusOK)	
		w.Write([]byte(fmt.Sprintf("%s/%s/%s", r.Host, "decode", shortUrl)))
	}else{
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Printf("unacceptable POST request")

		if len(url) < 1 {
			log.Println("body 'url' is missing")
			
			w.Write([]byte("body 'url' is missing"))
			return
		}
		w.Write([]byte("Request Header must be : application/x-www-form-urlencoded"))
		return
	}
}

func (h *HandlersLog) DecodeHandler(w http.ResponseWriter, r *http.Request){
	shortUrl := r.URL.Path[len("/decode/"):]

	if len(shortUrl) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("there is nothing to decode"))
		return
	}

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

	if len(shortUrl) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("there is nothing to redirect"))
		return
	}

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
		code := 200	
		startTime := time.Now()
		defer func() { 
			duration := time.Since(startTime)
			h.histogram.WithLabelValues(fmt.Sprintf("%d", code)).Observe(duration.Seconds())
			h.logger.Printf("request processed in %s\n", time.Now().Sub(startTime));
		}()
		next(w, r)
	}
}


func (h *HandlersLog) SetUpRoutes(mux *http.ServeMux){
	mux.HandleFunc("/encode", h.Logger(h.EncodeHandler))
	mux.HandleFunc("/decode/", h.Logger(h.DecodeHandler))
	mux.HandleFunc("/redirect/", h.Logger(h.RedirectHandler))
	mux.Handle("/metrics", prometheus.Handler())
}


func NewHandlersLog(logger *log.Logger, dbUrl string,  histogram *prometheus.HistogramVec) *HandlersLog{
	return &HandlersLog{
		logger: logger,
		dbUrl: dbUrl,
		histogram: histogram,
	}
}

func prometheusHandler(w http.ResponseWriter) http.Handler {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return prometheus.Handler()
}
