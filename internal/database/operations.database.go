package database

import (
	"database/sql"
	"example/web-service-gin/internal/repositories"
	"log"

	// Import the go-mssqldb driver anonymously
	_ "github.com/microsoft/go-mssqldb"
)

func ConnectToDB() *sql.DB {
	// Connection string in URL format
	// Example: "sqlserver://username:password@host/instance?param1=value&param2=value"
	connString := "sqlserver://sa:A7qmhn6vO9RxpRzwGE7AhR2ZkEfEPUHtOWBxuNaCydZGljv6CgfftIj6vfO@76.13.170.90?database=agendaHomologacao&encrypt=disable"

	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Error opening database connection: ", err.Error())
	}

	// It's a good practice to verify the connection is live
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err.Error())
	}
	log.Println("Successfully connected to SQL Server!")
	repositories.SearchClients(db)
	return db
}
