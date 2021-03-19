package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ProFL/gophercises-urlshort/models"
)

func JSONHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	redirectSpecs := make([]models.Redirect, 0)
	err := json.Unmarshal(data, &redirectSpecs)
	if err != nil {
		return nil, err
	}
	return MapHandler(mapFromRedirects(redirectSpecs), fallback), nil
}
