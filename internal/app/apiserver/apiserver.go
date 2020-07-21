package apiserver

import (
	"net/http"
)

func Start(config *Config/*, store sqlstore.Store*/) error {

	srv := newServer(/*store*/)
	srv.Configure()

	return  http.ListenAndServe(config.BindAddr, srv)
}
