package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/timmilesdw/backup-controller/pkg/config"
)

func StartServer(conf config.Spec) {
	app := fiber.New()

	api := app.Group("/api")
	v0 := api.Group("/v0")

	v0.Get("/cronjobs", func(c *fiber.Ctx) error {
		return c.JSON(conf.Backups)
	})
	v0.Get("/storages", func(c *fiber.Ctx) error {
		storages := conf.Storages
		for _, storage := range conf.Storages {
			storage.S3.AccessKey = "********"
			storage.S3.ClientSecret = "********"
		}
		return c.JSON(storages)
	})
	v0.Get("/databases", func(c *fiber.Ctx) error {
		databases := conf.Databases
		for _, database := range conf.Databases {
			database.Password = "********"
		}
		return c.JSON(databases)
	})
	// v0.Get("/backups")

	logrus.Fatal(app.Listen(":3000"))
}
