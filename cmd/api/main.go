package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/biacibengamukulu/stream-file-locally/internal/domain/filesystem"
	infraStorage "github.com/biacibengamukulu/stream-file-locally/internal/infrastructure/storage"
	"github.com/biacibengamukulu/stream-file-locally/internal/interfaces/rest/handlers"
	"github.com/biacibengamukulu/stream-file-locally/sharded/cassandra"
	"github.com/biacibengamukulu/stream-file-locally/sharded/config"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("No .env file loaded: " + err.Error())
	}

	cfg := config.Load("stream-file-locally")
	fileStorage, cleanup, err := buildStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if cleanup != nil {
		defer cleanup()
	}

	svcFileSystem := filesystem.NewFileService(cfg, fileStorage)

	app := fiber.New()
	app = config.FiberAdditionalConfig(app)
	app.Get(cfg.RoutePrefix+"/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": cfg.ServiceName})
	})
	registerDocs(app, cfg.RoutePrefix)
	api := app.Group(cfg.RoutePrefix).Group("/api/v1")

	fileHandler := handlers.NewFileSystemHandler(svcFileSystem)
	fileHandler.Register(api)

	log.Printf("[%s] listening on :%s storage=%s routePrefix=%s", cfg.ServiceName, cfg.HTTPPort, cfg.StorageDriver, cfg.RoutePrefix)
	log.Fatal(app.Listen(":" + cfg.HTTPPort))
}

func registerDocs(app *fiber.App, routePrefix string) {
	prefix := strings.TrimRight(routePrefix, "/")
	app.Get(prefix+"/openapi.yaml", func(c *fiber.Ctx) error {
		c.Type("yaml")
		return c.SendFile("./docs/openapi.yaml")
	})
	app.Get(prefix+"/swagger", func(c *fiber.Ctx) error {
		c.Type("html")
		return c.SendString(swaggerHTML(prefix + "/openapi.yaml"))
	})
}

func swaggerHTML(specURL string) string {
	return `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Stream File Locally API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <style>body{margin:0;background:#fff}</style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function() {
      window.ui = SwaggerUIBundle({
        url: "` + specURL + `",
        dom_id: "#swagger-ui",
        deepLinking: true,
        presets: [SwaggerUIBundle.presets.apis],
        layout: "BaseLayout"
      });
    };
  </script>
</body>
</html>`
}

func buildStorage(cfg config.Config) (filesystem.Storage, func(), error) {
	switch strings.ToLower(strings.TrimSpace(cfg.StorageDriver)) {
	case "", "disk", "filesystem", "file":
		storage, err := infraStorage.NewDiskStorage(cfg.DiskStoragePath)
		return storage, nil, err
	case "cassandra":
		rf, err := strconv.Atoi(cfg.CassandraReplicationFactor)
		if err != nil || rf <= 0 {
			rf = 1
		}
		if err := cassandra.BootstrapKeyspaceAndTables(cfg.CassandraHosts, cfg.CassandraKeyspace, rf, infraStorage.CassandraDDLs()); err != nil {
			return nil, nil, err
		}
		session, err := cassandra.New(cfg.CassandraHosts, cfg.CassandraKeyspace)
		if err != nil {
			return nil, nil, err
		}
		storage, err := infraStorage.NewCassandraStorage(session)
		if err != nil {
			session.Close()
			return nil, nil, err
		}
		return storage, session.Close, nil
	default:
		return nil, nil, fmt.Errorf("unsupported STORAGE_DRIVER %q", cfg.StorageDriver)
	}
}
