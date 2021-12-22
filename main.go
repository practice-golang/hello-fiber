package main // import "hello-fiber"

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/goccy/go-json"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Index(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}

func Hello(c *fiber.Ctx) error {
	name := c.Params("name")
	all := c.Params("*")
	log.Println(all)

	return c.SendString("Hello, " + name + " 👋!")
}

func JsonPOST(c *fiber.Ctx) error {
	p := Person{}

	json.Unmarshal(c.Body(), &p)

	c.SendStatus(http.StatusAccepted)
	return c.JSON(p)
}

func main() {
	listen := "127.0.0.1:2918"

	cfg := fiber.Config{
		AppName: "hello-fiber",
		// DisableStartupMessage: true,
	}

	b := fiber.New(cfg)

	b.Get("/", Index)
	b.Get("/hello/:name", Hello)
	b.Post("/json-post", JsonPOST)

	log.Fatal(b.Listen(listen))
}
