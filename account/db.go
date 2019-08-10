package account

import "database/sql"

type Datasource interface {
	QueryUserHashedPassword(username string) ([]byte, error)
}

type ProdDS struct {
	DB *sql.DB
}

var ProdDB ProdDS

func init() {
	var err error
	//DB, err = sql.Open("postgres", "postgres://testing:testing@localhost/testing?sslmode=disable")
	DB, err := sql.Open("postgres", "postgres://mitrakm_user:mitrakm@localhost/mitrakm?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}


}
