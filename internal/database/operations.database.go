package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/microsoft/go-mssqldb"
)

func ConnectToDB() *sql.DB {
	// Example: "sqlserver://username:password@host/instance?param1=value&param2=value"
	fmt.Println("Connecting to database:", os.Getenv("DB_SERVER"))
	connString := "sqlserver://sa:A7qmhn6vO9RxpRzwGE7AhR2ZkEfEPUHtOWBxuNaCydZGljv6CgfftIj6vfO@76.13.170.90?database=agendaHomologacao&encrypt=disable"
	//connString := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&encrypt=true&TrustServerCertificate=false",
	//os.Getenv("DB_USERNAME"),
	//os.Getenv("DB_PASSWORD"),
	//os.Getenv("DB_SERVER"),
	//os.Getenv("DB_PORT"),
	//os.Getenv("DB_DATABASE"),
	//)
	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Error opening database connection: ", err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err.Error())
	}
	log.Println("Successfully connected to SQL Server!")
	return db
}
