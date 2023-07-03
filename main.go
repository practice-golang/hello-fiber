package main // import "hello-fiber"

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"

	// "github.com/goccy/go-json"
	json "github.com/bytedance/sonic"
)

//go:embed static/html/*
var Content embed.FS

var staticPath = "static/html"

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

func exportEmbedStatic() error {
	exportPath := "."

	err := os.MkdirAll(exportPath, 0755)
	if err != nil {
		return err
	}

	err = fs.WalkDir(Content, "static", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			filePath := filepath.Join(exportPath, path)
			err := os.MkdirAll(filepath.Dir(filePath), 0755)
			if err != nil {
				return err
			}

			srcFile, err := Content.Open(path)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			dstFile, err := os.Create(filePath)
			if err != nil {
				return err
			}
			defer dstFile.Close()

			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func initServer(listen string) *fiber.App {
	cfg := fiber.Config{
		AppName:               "hello-fiber",
		DisableStartupMessage: true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	}

	b := fiber.New(cfg)

	b.Get("/api", Index)
	b.Get("/hello/:name", Hello)
	b.Post("/json-post", JsonPOST)

	b.Get("/user/:id?", func(c *fiber.Ctx) error {
		return c.SendString(c.Params("id"))
	})

	b.Get("/user/+", func(c *fiber.Ctx) error {
		return c.SendString(c.Params("+"))
	})

	b.Static("/", staticPath)

	return b
}

func main() {
	err := exportEmbedStatic()
	if err != nil {
		log.Fatal(err)
	}

	listen := "127.0.0.1:4416"

	b := initServer(listen)

	fmt.Println("Listening on " + listen)
	log.Fatal(b.Listen(listen))
}
