package handlers

import (
	"log"
	"net/http"

	"github.com/ProFL/gophercises-urlshort/models"
)

func mapFromRedirects(specs []models.Redirect) map[string]string {
	redirectMap := make(map[string]string)
	for _, spec := range specs {
		redirectMap[spec.Path] = spec.Url
	}
	return redirectMap
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if redirectUrl := pathsToUrls[req.URL.EscapedPath()]; redirectUrl != "" {
			log.Printf("Redirecting from %s to %s\n", req.URL.EscapedPath(), redirectUrl)
			res.Header().Add("Location", redirectUrl)
			res.WriteHeader(http.StatusFound)
			return
		}
		fallback.ServeHTTP(res, req)
	}
}
