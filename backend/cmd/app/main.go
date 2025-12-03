package main

import (
	"log"

	"github.com/stephenafamo/bob"
	"github.com/jimi024rion/todo-go/internal/infrastructure/database"
	"github.com/jimi024rion/todo-go/internal/infrastructure/repository"
	"github.com/jimi024rion/todo-go/internal/presentation/handler"
	"github.com/jimi024rion/todo-go/internal/presentation/router"
	"github.com/jimi024rion/todo-go/internal/usecase/todo"
)

func main() {
	// 1. Initialize database connection
	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()
	log.Println("Database connection established")

	// 2. Dependency Injection
	// Create instances of each layer, injecting dependencies from inner layers to outer layers.
	bobExec := bob.NewDB(db)
	todoRepo := repository.NewTodoRepository(bobExec)
	todoUsecase := todo.NewUsecase(todoRepo)
	todoHandler := handler.NewTodoHandler(todoUsecase)
	log.Println("Dependencies injected")

	// 3. Setup router
	r := router.NewRouter(todoHandler)
	log.Println("Router setup complete")

	// 4. Start server
	port := "8080"
	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
