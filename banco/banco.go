package banco

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //Drive de conexão com mySQL
)

//func
func Conect() (*sql.DB, error){

	stringConection := "golang:golang@/devbook?charset=utf8&parseTime=True&loc=Local"

	db, erro := sql.Open("mysql", stringConection)
	if erro != nil{
		
		return nil, erro
	}

	if db.Ping(); erro != nil{
		return nil, erro
	}

	return db, nil
}