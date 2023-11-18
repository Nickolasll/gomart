package main

import (
	"net/http"

	"github.com/Nickolasll/gomart/internal/presentation"
)

func main() {
	presentation.ParseFlags()
	mux, err := presentation.ChiFactory()
	if err != nil {
		panic(err)
	}
	err = http.ListenAndServe(*presentation.ServerEndpoint, mux)
	if err != nil {
		panic(err)
	}
}
