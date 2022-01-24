package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/timmilesdw/backup-controller/pkg/backupper"
	"github.com/timmilesdw/backup-controller/pkg/config"
)

func StartServer(conf config.UI) {
	app := fiber.New()

	api := app.Group("/api")
	v0 := api.Group("/v0")

	// v0.Get("/cronjobs", func(c *fiber.Ctx) error {
	// 	return c.JSON()
	// })
	v0.Get("/storers", func(c *fiber.Ctx) error {
		storers := []map[string]interface{}{}
		for _, st := range backupper.Storers {
			storers = append(storers, st.GetMap())
		}

		return c.JSON(storers)
	})
	v0.Get("/exporters", func(c *fiber.Ctx) error {
		exps := []map[string]interface{}{}
		for _, ex := range backupper.Exporters {
			exps = append(exps, ex.GetMap())
		}

		return c.JSON(exps)
	})
	// v0.Get("/backups")

	logrus.Fatal(app.Listen(":3000"))
}
