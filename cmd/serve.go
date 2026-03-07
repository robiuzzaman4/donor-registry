package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/robiuzzaman4/donor-registry-backend/internal/config"
	"github.com/robiuzzaman4/donor-registry-backend/internal/infrastructure/db"
)

func Serve() {
	// config
	cnf := config.GetConfig()

	// context
	ctx := context.Background()

	// db conn
	dbConn, err := db.NewConnection(ctx, cnf.DbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Database connected:", dbConn)

}
