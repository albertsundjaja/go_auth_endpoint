package account

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func HandleUsers_GET(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")


}

