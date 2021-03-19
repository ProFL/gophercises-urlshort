package urlshort

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ProFL/gophercises-urlshort/models"
	yaml "gopkg.in/yaml.v2"
)

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

func mapFromRedirectSpecs(specs []models.Redirect) map[string]string {
	redirectMap := make(map[string]string)
	for _, spec := range specs {
		redirectMap[spec.Path] = spec.Url
	}
	return redirectMap
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	redirectSpecs := make([]models.Redirect, 0)
	err := yaml.Unmarshal(yml, &redirectSpecs)
	if err != nil {
		return nil, err
	}
	return MapHandler(mapFromRedirectSpecs(redirectSpecs), fallback), nil
}

func JSONHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	redirectSpecs := make([]models.Redirect, 0)
	err := json.Unmarshal(data, &redirectSpecs)
	if err != nil {
		return nil, err
	}
	return MapHandler(mapFromRedirectSpecs(redirectSpecs), fallback), nil
}
