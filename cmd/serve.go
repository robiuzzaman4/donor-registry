package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/robiuzzaman4/donor-registry-backend/internal/config"
	"github.com/robiuzzaman4/donor-registry-backend/internal/infrastructure/db"
	"github.com/robiuzzaman4/donor-registry-backend/internal/rest"
)

func Serve() {
	// config
	cnf := config.GetConfig()

	// context
	ctx := context.Background()

	// db connection
	dbConn, err := db.NewConnection(ctx, cnf.DbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// start server
	server := rest.NewServer(cnf, ctx, dbConn)
	server.Start()
}
