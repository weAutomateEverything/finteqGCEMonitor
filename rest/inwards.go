package rest

import (
	"net/http"
	"github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/service"
)

func initializeInwards(w http.ResponseWriter, r *http.Request){
	service.ParseInwardCutttoffTimes(r.Body)

}