package main

import (
	"github.com/julienschmidt/httprouter"
	//"github.com/rs/cors"
	account "go_admin/account"
	"log"
	"net/http"
)

func main() {
	router := httprouter.New()

	router.POST("/api/account/register", account.HandleRegister_POST(account.ProdDB))
	router.OPTIONS("/api/account/register", HandleOptions)
	router.POST("/api/account/login", account.HandleLogin_POST(account.ProdDB))
	router.OPTIONS("/api/account/login", HandleOptions)
	router.GET("/api/account/protected_page", account.HandleProtectedEndpoint_GET())
	router.OPTIONS("/api/account/protected_page", HandleOptions)

	//handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8000", router))
}

func HandleOptions(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, jwt, Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	w.WriteHeader(http.StatusOK)
}
