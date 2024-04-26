package main

import (
	"construction-organization-report/internal/api"
	"construction-organization-report/internal/db"
	"construction-organization-report/internal/log"
	"construction-organization-report/pkg/config"
	"database/sql"
)

func main() {
	dbConfig, err := config.NewDBConfig()
	if err != nil {
		log.Logger.WithError(err).Errorln("Db config isn't init")
		return
	}

	log.Logger.Infoln("Connect to database...")
	newDB, err := db.NewDB(dbConfig)
	if err != nil {
		log.Logger.WithError(err).Errorln("Error on database connection")
		return
	}

	defer func(newDB *sql.DB) {
		err := newDB.Close()
		if err != nil {
			log.Logger.WithError(err).Errorln("Error on database close")
		}
	}(newDB)

	log.Logger.Infoln("Successful database connection")

	log.Logger.Infoln("Run server on port 8080")
	server := api.NewServer(":8080", newDB)
	server.InitializeRoutes()
	server.Start()
}
