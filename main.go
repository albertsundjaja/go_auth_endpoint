package main

import (
	"github.com/julienschmidt/httprouter"
	account "go_admin/account"
	"log"
	"net/http"
)

func main() {
	router := httprouter.New()

	router.POST("/api/account/register", account.HandleRegister_POST(account.ProdDB))
	router.POST("/api/account/login", account.HandleLogin_POST(account.ProdDB))
	router.GET("/api/account/protected_page", account.HandleProtectedEndpoint_GET())

	log.Fatal(http.ListenAndServe(":8000", router))
}
