package base

import(
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

)

func ConnectToDb()(*sql.DB,error){
	databaseHost :="localhost"
	databaseUser :="student"
	dbPassword := base.GetEnv("DATABASE_PASSWORD", "") 
	if dbPassword == "" {
		log.Fatal("DATABASE_PASSWORD environment variable is not set")
	}
	connection := fmt.Sprintf("postgresql://student:%s@%s:5432/proxy_db?sslmode=disable", dbPassword, databaseHost)

	var db *sql.DB
	var err error
	db, err = sql.Open("postgres", connection)
	if err !=nil{
		return nil, err
	}	
	return db,nil
}

