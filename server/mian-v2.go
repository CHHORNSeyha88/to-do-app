package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
	Body  string `json:"body"`
}

var (
	todos []Todo
	mu    sync.RWMutex
)

func main() {
	app := fiber.New()

	// Middlewares
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: getEnv("CORS_ORIGIN", "http://localhost:5173"),
		AllowMethods: "GET,POST,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Healthcheck
	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("Ok")
	})

	api := app.Group("/api/v1")
	todosAPI := api.Group("/todos")

	// Create todo
	todosAPI.Post("/", func(c *fiber.Ctx) error {
		var in struct {
			Title string `json:"title"`
			Body  string `json:"body"`
		}
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
		}
		if strings.TrimSpace(in.Title) == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title is required"})
		}

		mu.Lock()
		todo := Todo{ID: len(todos) + 1, Title: in.Title, Body: in.Body, Done: false}
		todos = append(todos, todo)
		mu.Unlock()

		return c.Status(fiber.StatusCreated).JSON(todo)
	})

	// Mark as done
	todosAPI.Patch("/:id/done", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		mu.Lock()
		defer mu.Unlock()
		for i := range todos {
			if todos[i].ID == id {
				todos[i].Done = true
				return c.JSON(todos[i])
			}
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "todo not found"})
	})

	// List todos
	todosAPI.Get("/", func(c *fiber.Ctx) error {
		mu.RLock()
		defer mu.RUnlock()
		return c.JSON(todos)
	})

	// Get by id
	todosAPI.Get("/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		mu.RLock()
		defer mu.RUnlock()
		for _, t := range todos {
			if t.ID == id {
				return c.JSON(t)
			}
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "todo not found"})
	})

	// Start server in a goroutine to enable graceful shutdown
	go func() {
		if err := app.Listen(":4000"); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")
	_ = app.Shutdown()
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}


