package main

import (
	"github.com/julienschmidt/httprouter"
	account "go_admin/account"
)

func main() {
	router := httprouter.New()

	router.GET("/api/admin/users", account.HandleUsers_GET)

}
