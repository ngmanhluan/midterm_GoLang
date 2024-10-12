package gin

import (
	"net/http"
	"todo-app/domain"
	"todo-app/pkg/clients"
	"todo-app/pkg/tokenprovider"
	"github.com/google/uuid" 
	"github.com/gin-gonic/gin"
)

type UserService interface {
	Register(data *domain.UserCreate) error
	Login(data *domain.UserLogin) (tokenprovider.Token, error)
	GetAllUsers() ([]*domain.User, error)
	GetUserByID(id uuid.UUID) (*domain.User, error)
	UpdateUser(id uuid.UUID, firstName string, lastName string) error
	DeleteUser(id uuid.UUID) error
}

type userHandler struct {
	userService UserService
}

func NewUserHandler(apiVersion *gin.RouterGroup, svc UserService) {
	userHandler := &userHandler{
		userService: svc,
	}

	users := apiVersion.Group("/users")
	users.POST("/register", userHandler.RegisterUserHandler)
	users.POST("/login", userHandler.LoginHandler)
	users.GET("/", userHandler.GetAllUsersHandler)
	users.GET("/:id", userHandler.GetUserByIDHandler)
	users.PUT("/:id", userHandler.UpdateUserHandler)
	users.DELETE("/:id", userHandler.DeleteUserHandler)
}

func (h *userHandler) RegisterUserHandler(c *gin.Context) {
	var data domain.UserCreate

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))
		return
	}

	if err := h.userService.Register(&data); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return	
	}	

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse(data.ID))
}


func (h *userHandler) LoginHandler(c *gin.Context) {
	var data domain.UserLogin

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))
		return
	}

	token, err := h.userService.Login(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse(token))
}

func (h *userHandler) GetAllUsersHandler(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *userHandler) GetUserByIDHandler(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id) 
	if err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandler) UpdateUserHandler(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id) 
	if err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))
		return
	}

	var data struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))
		return
	}

	if err := h.userService.UpdateUser(userID, data.FirstName, data.LastName); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse("User updated successfully"))
}

func (h *userHandler) DeleteUserHandler(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id) 
	if err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))
		return
	}

	if err := h.userService.DeleteUser(userID); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse("User deleted successfully"))
}
