package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/biangacila/luvungula-go/global"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func FiberAdditionalConfig(app *fiber.App) *fiber.App {
	app = fiber.New(fiber.Config{
		ReadTimeout:  0,
		WriteTimeout: 0,
		BodyLimit:    50 * 1024 * 1024, // 50MB to accommodate base64 images
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			if code == fiber.StatusNotFound {
				log.Printf("404 NOT FOUND: %s %s from %s",
					c.Method(),
					c.OriginalURL(),
					c.IP(),
				)
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Recover from panics so the process doesn't die
	app.Use(recover.New())

	// Logger: skip the SSE stream path to reduce log spam
	app.Use(logger.New(logger.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/backend-swalylogisticsdriver/api/notifications/stream"
		},
	}))

	// --- SAFE MULTI-ORIGIN CORS SETUP ---
	rawOrigins := os.Getenv("CORS_ORIGINS")
	var origins []string
	for _, o := range strings.Split(rawOrigins, ",") {
		o = strings.TrimSpace(o)
		if o != "" && o != "*" {
			origins = append(origins, o)
		}
	}
	if len(origins) == 0 {
		log.Println("⚠️  No valid CORS origins found, using localhost as fallback")
		origins = []string{"http://localhost:5173"}
	}

	subdomains := []string{
		".biacibenga.com",
		".webcontainer-api.io",
		".easipath.com",
		".webcontainer-api.io",
	}
	mainDomains := []string{
		"https://biacibenga.com",
		"webcontainer-api.io",
		"easipath.com",
		"https://zp1v56uxy8rdx5ypatb0ockcb9tr6a-oci3--5173--61636aac.local-credentialless.webcontainer-api.io",
	}

	global.DisplayObject("All CORS", origins)

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			fmt.Println("------:) CORS: ", origin)
			// allow localhost for development
			if strings.Contains(origin, "localhost") {
				return true
			}

			// allow all subdomains of biacibenga.com
			for _, subdomain := range subdomains {
				if strings.HasSuffix(origin, subdomain) {
					return true
				}
			}

			// allow main domain
			for _, domain := range mainDomains {
				if origin == domain {
					return true
				}
			}

			return false
		},
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	return app
}
