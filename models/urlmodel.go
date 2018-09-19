package models

import (
	"fmt"

	"github.com/OtgonbayarT/microservice/controllers"
	"github.com/rapidloop/skv"
)

func InsertUrl(dbUrl string, longUrl string) (string, error) {
	store, err := skv.Open(dbUrl)
	if err != nil {		
		return "", err
	}
	defer store.Close()

	shortUrl := fmt.Sprint(controllers.Hash(longUrl))

	if err := store.Put(shortUrl, longUrl); err != nil {
		return "", err
	}

	return shortUrl, nil
}

func GetUrl(dbUrl string, shortUrl string) (string, error) {
	store, err := skv.Open(dbUrl)
	if err != nil {		
		return "", err
	}
	defer store.Close()

	var val string
	if err := store.Get(fmt.Sprint(shortUrl), &val); err != nil {		
		return "", err
	}	

	if len(val) == 0 {
		return "there is LongUrl for matching shortUrl!!!", nil
	}
	return val, nil
}