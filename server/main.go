package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
	Body  string `json:"body"`
}

func main() {
	fmt.Println("HelloWorld")
	// declear variable
	app := fiber.New()
	todos := []Todo{}

	// get health-check
	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("Ok")
	})

	// insert
	app.Post("api/v1/todos", func(c *fiber.Ctx) error {

		todo := &Todo{}
		if err := c.BodyParser(todo); err != nil {
			return err
		}
		todo.ID = len(todos) + 1

		todos = append(todos, *todo)
		return c.JSON(todos)

	})

	// update
	app.Patch("api/v1/todos/:id/done", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")

		if err != nil {
			return c.Status(401).SendString("invalid id")
		}
		for i, t := range todos {
			if t.ID == id {
				todos[i].Done = true
				break
			}
		}

		return c.JSON(todos)

	})

	// get

	app.Get("api/v1/todos", func(c *fiber.Ctx) error {
		return c.JSON(todos)
	})

	// get by id

	app.Get("api/v1/todos/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")

		if err != nil {
			return c.Status(401).SendString("invalid id")
		}

		// 일치하는 ID가 있는 항목을 찾기 위해 todos를 반복합니다.

		for _, t := range todos {
			if t.ID == id {
				// Found the todo — return it as JSON
				return c.JSON(t)
			}
		}
		return c.Status(400).SendString("Todo not found!!")

	})

	log.Fatal(app.Listen(":4000"))
}
