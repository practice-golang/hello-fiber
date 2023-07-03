package main // import "hello-fiber"

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"

	// "github.com/goccy/go-json"
	json "github.com/bytedance/sonic"
)

//go:embed html/*
var Content embed.FS

var staticPath = "../html"

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
	listen := "127.0.0.1:4416"

	cfg := fiber.Config{
		AppName:               "hello-fiber",
		DisableStartupMessage: true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	}

	b := fiber.New(cfg)

	b.Get("/", Index)
	b.Get("/hello/:name", Hello)
	b.Post("/json-post", JsonPOST)

	b.Get("/user/:id?", func(c *fiber.Ctx) error {
		return c.SendString(c.Params("id"))
	})

	b.Get("/user/+", func(c *fiber.Ctx) error {
		return c.SendString(c.Params("+"))
	})

	b.Use("/*.html", filesystem.New(filesystem.Config{
		Root:       http.FS(Content),
		PathPrefix: "html",
		Browse:     true,
	}))

	b.Static("/", staticPath)
	// b.Use("/", filesystem.New(filesystem.Config{
	// 	Root:       http.FS(Content),
	// 	PathPrefix: "html",
	// 	Browse:     true,
	// }))

	fmt.Println("Listening on " + listen)
	log.Fatal(b.Listen(listen))
}
