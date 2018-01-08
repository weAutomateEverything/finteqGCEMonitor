package rest

import (
	"net/http"
	"github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/service"
)

func initializeInwards(w http.ResponseWriter, r *http.Request){
	service.ParseInwardCutttoffTimes(r.Body)
}

func pingHal(w http.ResponseWriter, r *http.Request){
	err := service.Ping()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}