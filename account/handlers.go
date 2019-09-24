package account

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/syndtr/goleveldb/leveldb"
	"net/http"
	"time"
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
		var func_name = "HandleLogin_POST"

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, jwt, Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
		w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS")

		var loginReq LoginRequest
		err := ReadRequestBody(req, &loginReq)
		if err != nil {
			fmt.Println(func_name, err)
			HttpErrorResponder(w, "Error decoding req body")
			return
		}

		registeredPass, err := ds.QueryUserPassword(loginReq.Username)
		if err != nil {
			if err == leveldb.ErrNotFound {
				fmt.Println(func_name, "wrong username")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("Wrong Username/Password"))
				return
			}
			fmt.Println(func_name, err)
			HttpErrorResponder(w, "DB error")
			return
		}

		if registeredPass != loginReq.Password {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Wrong Username/Password"))
			return
		}

		jwt, err := CreateJwt(loginReq.Username, time.Time{})
		if err != nil {
			fmt.Println(func_name, err)
			HttpErrorResponder(w, "JWT error")
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(jwt))
	})
}

func HandleRegister_POST(ds Datasource) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		var func_name = "HandleRegister_POST"

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, jwt, Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
		w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS")

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

func HandleProtectedEndpoint_GET() httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "jwt, UUID, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

		var func_name = "HandleProtectedEndpoint"

		jwt := req.Header.Get("jwt")
		if jwt == "" {
			fmt.Println(func_name, "cannot get jwt header")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprint(w, "Jwt header invalid")
			return
		}

		claims, err := ValidateSigningAndGetJwtClaims(jwt)
		if err != nil {
			fmt.Println(func_name, err)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprint(w, "Jwt invalid or expired")
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello " + claims["Username"].(string)))
	})
}
