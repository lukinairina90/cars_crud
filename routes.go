package main

import (
	"C_CRUD/car"
	"github.com/go-chi/chi/v5"
)

func initRouter(carTransport *car.Transport) *chi.Mux {
	router := chi.NewRouter()

	router.Route("/api", func(r chi.Router) {
		r.Route("/cars", func(r chi.Router) {
			r.Get("/", carTransport.ShowCars)
			r.Post("/", carTransport.CreateCar)
			r.Put("/{id:[0-9]+}", carTransport.EditCarById)
			r.Get("/{id:[0-9]+}", carTransport.ShowCarByID)
			r.Delete("/{id:[0-9]+}", carTransport.DeleteCar)
		})
	})

	router.Get("/", carTransport.ShowStaticSinglePageApplication)

	return router
}
