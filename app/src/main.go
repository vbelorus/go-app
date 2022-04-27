package main

import (
	"github.com/ClickHouse/clickhouse-go/v2"
	clickhouse_model "github.com/vbelorus/go-app/v2/src/clickhouse"
	"github.com/vbelorus/go-app/v2/src/config"
	"github.com/vbelorus/go-app/v2/src/models"
	"github.com/vbelorus/go-app/v2/src/router"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	httpServerPort := LookupEnvWithDefaultValue("HTTP_SERVER_PORT", "8081")
	clickhouseHost := LookupEnvWithDefaultValue("CLICKHOUSE_HOST", "clickhouse-server")
	clickhousePort := LookupEnvWithDefaultValue("CLICKHOUSE_PORT", "9000")
	clickhouseDatabase := LookupEnvWithDefaultValue("CLICKHOUSE_DATABASE", "app")
	clickhouseUsername := LookupEnvWithDefaultValue("CLICKHOUSE_USERNAME", "default")
	clickhousePassword := LookupEnvWithDefaultValue("CLICKHOUSE_PASSWORD", "")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	deviceEventChannel := make(chan models.DeviceEvent)
	errorsChannel := make(chan error)

	var (
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{clickhouseHost + ":" + clickhousePort},
			Auth: clickhouse.Auth{
				Database: clickhouseDatabase,
				Username: clickhouseUsername,
				Password: clickhousePassword,
			},
			//Debug:           true,
			DialTimeout:     time.Second,
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: time.Hour,
		})
	)

	if err != nil {
		log.Fatal(err)
	}

	app := &config.Application{
		InfoLog:            infoLog,
		ErrorLog:           errorLog,
		DeviceEventChannel: deviceEventChannel,
		AppDB:              &clickhouse_model.AppDB{Conn: conn, ErrorsChannel: errorsChannel},
	}

	srv := &http.Server{
		Addr:    ":" + httpServerPort,
		Handler: router.SetRoutes(app),
	}

	go app.AppDB.ListenDeviceEvents(deviceEventChannel)
	go func() {
		for errorElement := range errorsChannel {
			app.LogError(errorElement)
		}
	}()
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func LookupEnvWithDefaultValue(key string, defaultValue string) string {
	envValue, ok := os.LookupEnv(key)
	if !ok {
		envValue = defaultValue
	}

	return envValue
}
