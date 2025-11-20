package main

import (
	"net/http"
)

func routes(handler *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /links", handler.CheckLinksHandler)
	mux.HandleFunc("GET /report", handler.GenerateReportHandler)
	return mux
}
