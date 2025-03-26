package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/victor-tsykanov/delivery/cmd/app"
	"github.com/victor-tsykanov/delivery/cmd/http"
	"github.com/victor-tsykanov/delivery/cmd/jobs"
	"github.com/victor-tsykanov/delivery/internal/common/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config.MustLoadEnv(".env")
	dbConfig := config.MustLoadDBConfig()
	httpConfig := config.MustLoadHTTPConfig()
	geoServiceConfig := config.MustLoadGeoServiceConfig()

	ctx := context.Background()
	db := mustConnectToDB(dbConfig.DSN())
	root := app.NewCompositionRoot(ctx, db, geoServiceConfig.Address)

	go http.Serve(ctx, root, httpConfig)
	go jobs.AssignOrders(ctx, root)
	go jobs.MoveCouriers(ctx, root)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	<-signalChan

	root.Shutdown(ctx)
}

func mustConnectToDB(dsn string) *gorm.DB {
	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	return db
}
