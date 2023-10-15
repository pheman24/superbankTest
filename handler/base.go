package handler

import "net/http"

func New() http.Handler {
	mux := http.NewServeMux()
	// Root
	mux.Handle("/",  http.FileServer(http.Dir("template/")))

	// OauthGoogle
	mux.HandleFunc("/auth/google/login", oauthGoogleLogin)
	mux.HandleFunc("/auth/google/callback", oauthGoogleCallback)
	mux.HandleFunc("/auth/google/loginByEmail", LoginByEmail)
	mux.HandleFunc("/auth/google/submitForm", SubmitForm)

	return mux
}