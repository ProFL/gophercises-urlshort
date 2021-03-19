package handlers

import (
	"net/http"

	"github.com/ProFL/gophercises-urlshort/models"
	yaml "gopkg.in/yaml.v2"
)

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
	return MapHandler(mapFromRedirects(redirectSpecs), fallback), nil
}
