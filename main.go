package main // import "hello-fiber"

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	// "github.com/goccy/go-json"
	// json "github.com/bytedance/sonic"
)

//go:embed static/html/*
var Content embed.FS

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var UploadRoot = "upload"

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
	jsonData := Person{}

	// err := c.BodyParser(&jsonData)
	err := json.Unmarshal([]byte(c.FormValue("person")), &jsonData)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	fdata, err := c.FormFile("image")
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusBadRequest).SendString("Bad Request")
	}

	fmt.Println(fdata.Filename)
	c.SaveFile(fdata, UploadRoot+"/"+fdata.Filename)

	return c.Status(http.StatusOK).JSON(jsonData)
}

func initServer(listen string) *fiber.App {
	if _, err := os.Stat(UploadRoot); os.IsNotExist(err) {
		os.Mkdir(UploadRoot, os.ModePerm)
	}

	// engine := html.New("./static/html", ".html")
	engine := html.NewFileSystem(http.FS(Content), ".html")
	cfg := fiber.Config{
		AppName:               "hello-fiber",
		DisableStartupMessage: true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		Views:                 engine,
	}

	embedRoot := "static/html"

	app := fiber.New(cfg)

	app.Get("/api", Index)
	app.Get("/hello/:name", Hello)
	app.Post("/json-post", JsonPOST)

	app.Get("/user/:id?", func(c *fiber.Ctx) error {
		return c.SendString(c.Params("id"))
	})

	app.Get("/user/+", func(c *fiber.Ctx) error {
		return c.SendString(c.Params("+"))
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render(embedRoot+"/index", fiber.Map{})
	})

	app.Get("/hello-tmpl", func(c *fiber.Ctx) error {
		return c.Render(embedRoot+"/template_admin/index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	return app
}

func main() {
	listen := "127.0.0.1:4416"

	b := initServer(listen)

	fmt.Println("Listening on " + listen)
	log.Fatal(b.Listen(listen))
}
