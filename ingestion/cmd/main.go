package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/config"
)

const configPath = "config.yml"

func main() {
	config.MustLoad(configPath)

	a := app.New()

	if err := a.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "app run failed: %v\n", err)
		os.Exit(1)
	}
}
