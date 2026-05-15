package main

import (
	"fmt"
	"log"

	"github.com/biacibengamukulu/stream-file-locally/internal/domain/filesystem"
	"github.com/biacibengamukulu/stream-file-locally/internal/interfaces/rest/handlers"
	"github.com/biacibengamukulu/stream-file-locally/sharded/config"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	v := config.Getenv("CASSANDRA_HOSTS", "")
	if v == "" {
		if err := godotenv.Load(".env"); err != nil {
			fmt.Println("Error loading .env file: " + err.Error())
		}

		cfg := config.Load("ledger")
		hosts := config.Getenv("CASSANDRA_HOSTS", cfg.CassandraHosts)
		keyspace := config.Getenv("CASSANDRA_KEYSPACE", cfg.CassandraKeyspace)
		routePrefix := config.Getenv("ROUTE_PREFIX", "/biatechwallet-ledger")
		port := config.Getenv("HTTP_PORT", "18203")

		fmt.Println(hosts, keyspace, routePrefix, port)

		//todo repositories  casssandra here

		// services
		svcFileSystem := filesystem.NewFileService(cfg)

		app := fiber.New()
		app = config.FiberAdditionalConfig(app)
		app.Get(routePrefix+"/health", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ok"}) })
		api := app.Group(routePrefix).Group("/api/v1")

		fileHandler := handlers.NewFileSystemHandler(svcFileSystem)
		fileHandler.Register(api)

		log.Printf("[ledger] listening on :%s (keyspace=%s) routePrefix=%s", port, keyspace, routePrefix)
		log.Fatal(app.Listen(":" + port))

	}

}
