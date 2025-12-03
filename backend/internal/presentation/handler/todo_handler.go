package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/internal/usecase/todo"
)

type TodoHandler struct {
	todoUsecase todo.Usecase
}

func NewTodoHandler(tu todo.Usecase) *TodoHandler {
	return &TodoHandler{todoUsecase: tu}
}

type createTodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// Create handles POST /todos requests.
func (h *TodoHandler) Create(c *gin.Context) {
	var req createTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	createdTodo, err := h.todoUsecase.CreateTodo(c.Request.Context(), req.Title, req.Description)
	if err != nil {
		// A more sophisticated error handling would be better.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create todo"})
		return
	}

	c.JSON(http.StatusCreated, createdTodo)
}

// GetByID handles GET /todos/:id requests.
func (h *TodoHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	t, err := h.todoUsecase.GetTodoByID(c.Request.Context(), id)
	if err != nil {
		if err == todo.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get todo"})
		return
	}

	c.JSON(http.StatusOK, t)
}

// GetAll handles GET /todos requests.
func (h *TodoHandler) GetAll(c *gin.Context) {
	todos, err := h.todoUsecase.GetAllTodos(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get todos"})
		return
	}

	c.JSON(http.StatusOK, todos)
}

type updateTodoRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

// Update handles PUT /todos/:id requests.
func (h *TodoHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	var req updateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	updatedTodo, err := h.todoUsecase.UpdateTodo(c.Request.Context(), id, req.Title, req.Description, req.Completed)
	if err != nil {
		if err == todo.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update todo"})
		return
	}

	c.JSON(http.StatusOK, updatedTodo)
}

// Delete handles DELETE /todos/:id requests.
func (h *TodoHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	err = h.todoUsecase.DeleteTodo(c.Request.Context(), id)
	if err != nil {
		if err == todo.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete todo"})
		return
	}

	c.Status(http.StatusNoContent)
}
