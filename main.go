package main

import (
	"github.com/julienschmidt/httprouter"
	account "go_admin/account"
	"net/http"
)

func main() {
	router := httprouter.New()

	//router.GET("/api/admin/users", account.HandleUsers_GET)
	router.POST("/api/account/register", account.HandleregisterPost(account.ProdDB))

	http.ListenAndServe(":8000", router)
}
