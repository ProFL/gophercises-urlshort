package handlers

import (
	"log"
	"net/http"

	"github.com/ProFL/gophercises-urlshort/repositories"
)

func DatabaseHandler(repo *repositories.RedirectRepository, fallback http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		spec, err := repo.FindByPath(req.URL.EscapedPath())
		if spec == nil || err != nil {
			if err != nil {
				log.Println("Failed to fetch redirect from the database", err.Error())
			}
			fallback.ServeHTTP(res, req)
			return
		}
		log.Println("Redirecting from", spec.Path, "to", spec.Url)
		res.Header().Add("Location", spec.Url)
		res.WriteHeader(http.StatusFound)
	}
}
