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
			msg := "body 'url' is missing"
			log.Println(msg)
			
			w.Write([]byte(msg))
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

		startTime := time.Now()
		o := &responseObserver{ResponseWriter: w}
		defer func() { 			
			duration := time.Since(startTime)		    
			h.logger.Printf("http status %d\n", o.status);
			h.histogram.WithLabelValues(fmt.Sprintf("%d", o.status)).Observe(duration.Seconds())
			h.logger.Printf("request processed in %s\n", time.Now().Sub(startTime));
		}()
		next.ServeHTTP(o, r)
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

type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (o *responseObserver) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

func (o *responseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}