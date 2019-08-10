package account

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// ReadRequestBody will read the incoming request body, decode it into json and write into the inputted struct pointer
func ReadRequestBody(req *http.Request, result interface{}) (err error) {
	var func_name = "ReadRequestBody"

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&result)
	if err != nil {
		fmt.Println(func_name, "Error decoding request body, check if you pass the correct interface", err)
		return err
	}

	return nil
}

// A method to response internal server error with given message
func HttpErrorResponder(w http.ResponseWriter, msg string) {
	var func_name = "HttpErrorResponder"
	w.WriteHeader(http.StatusInternalServerError)
	_,err := w.Write([]byte(msg))
	if err != nil {
		fmt.Println(func_name, err)
	}
}

func HandleLogin_POST(ds Datasource) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		var func_name = "HandleLogin_GET"

		var loginReq LoginRequest
		err := ReadRequestBody(req, &loginReq)
		if err != nil {
			fmt.Println(func_name, err)
			HttpErrorResponder(w, "Error decoding req body")
			return
		}

		registeredPass, err := ds.QueryUserHashedPassword(loginReq.Username)
		if err != nil {
			fmt.Println(func_name, err)
			HttpErrorResponder(w, "DB error")
			return
		}

		if registeredPass != loginReq.Password {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Wrong Username/Password"))
			return
		}



	})
}

func HandleregisterPost(ds Datasource) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		var func_name = "HandleRegister_POST"

		var registerReq RegisterRequest
		err := ReadRequestBody(req, &registerReq)
		if err != nil {
			fmt.Println(func_name, err)
			HttpErrorResponder(w, "error decoding req body")
			return
		}

		err = ds.InputNewUser(registerReq)
		if err != nil {
			fmt.Println(func_name, err)
			HttpErrorResponder(w, "DB error")
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("Account Created"))
	})
}