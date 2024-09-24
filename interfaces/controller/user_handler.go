package controller

import (
	"fmt"
	"food-app/application"
	"food-app/domain/entity"
	"food-app/infrastructure/auth"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Users struct defines the dependencies that will be used
type Users struct {
	us application.UserAppInterface
	rd auth.AuthInterface
	tk auth.TokenInterface
}

// Users constructor
func NewUsers(us application.UserAppInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Users {
	return &Users{
		us: us,
		rd: rd,
		tk: tk,
	}
}

func (s *Users) SaveUser(c *gin.Context) {
	var user entity.User
	// In  first name and last name

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"invalid_json": "invalid json",
		})
		return
	}
	fmt.Println("FirstName: ", user.FirstName)
	fmt.Println("LastName: ", user.LastName)
	validateErr := user.Validate("")
	if len(validateErr) > 0 {
		c.JSON(http.StatusBadRequest, validateErr)
		return
	}
	newUser, err := s.us.SaveUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, newUser.PublicUser())
}

func (s *Users) GetUsers(c *gin.Context) {
	users := entity.Users{} //customize user
	var err error
	//us, err = application.UserApp.GetUsers()
	users, err = s.us.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, users.PublicUsers())
}

func (s *Users) GetUser(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user, err := s.us.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user.PublicUser())
}
