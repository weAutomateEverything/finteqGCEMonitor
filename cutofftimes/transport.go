package cutofftimes

import (
	"github.com/zamedic/go2hal/gokit"
	"github.com/gorilla/mux"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"

	"net/http"
)

func MakeHandler(service Service, logger kitlog.Logger) http.Handler {
	opts := gokit.GetServerOpts(logger)


	alertHandler := kithttp.NewServer(makeCutoffTimesEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts..., )


	r := mux.NewRouter()

	r.Handle("/inwards/cutoff", alertHandler).Methods("POST")

	return r
}
