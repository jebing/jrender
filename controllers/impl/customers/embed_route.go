package customers

import "github.com/go-chi/chi/v5"

func Route(r chi.Router) {

	r.Route("/api/v1/customers", func(r chi.Router) {
	})
}
