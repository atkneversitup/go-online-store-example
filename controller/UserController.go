package controller

import (
	"fmt"
	"log"
	"myapp-go-echo/database"
	"os"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

// "myapp-go-echo/database"

// UserController struct
type UserController struct{}
type Login struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}

// GetUsers function
func (uc *UserController) GetUsers(c echo.Context) error {
	var users []database.User
	if err := database.DB.Find(&users).Error; err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, users)
}
func (uc *UserController) GetUsersById(c echo.Context) error {
	var user database.User
	if err := database.DB.Find(&user).Error; err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, user)
}
func (uc *UserController) CreateUsers(c echo.Context) error {
	var user database.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(400, err)
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, user)
}

// Register user function
func (uc *UserController) Register(c echo.Context) error {
	// Hash the password
	var user database.User
	log.Println(`User: `, c)
	if err := c.Bind(&user); err != nil {
		return c.JSON(400, err)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Error hashing password: %v", err)
	}
	// Save the hashed password to the user struct
	user.Password = string(hashedPassword)

	// Save the user to the database
	// ...
	if err := database.DB.Create(&user).Error; err != nil {
		return c.JSON(500, err)
	}
	log.Println(`User created: `, user)
	return c.JSON(200, user)
}

// Login user function
func (uc *UserController) Login(c echo.Context) error {
	var login Login
	if err := c.Bind(&login); err != nil {
		return c.JSON(400, err)
	}
	var user database.User
	if err := database.DB.Where("username = ?", login.Username).First(&user).Error; err != nil {
		return c.JSON(500, err)
	}
	// Compare the provided password with the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		log.Println(`Error comparing passwords: `, err)
		return c.JSON(401, map[string]string{"error": err.Error()})
	}
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Username
	claims["admin"] = true
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return err
	}
	return c.JSON(200, map[string]string{
		"token": t,
	})
}
func (uc *UserController) ComparePassword(user *database.User, password string) error {
	// Compare the provided password with the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return fmt.Errorf("Error comparing passwords: %v", err)
	}
	return nil
}

// func (uc *UserController) createUsers(c echo.Context) error {
// 	var user database.User
// 	if err := c.Bind(&user); err != nil {
// 		return c.JSON(400, err)
// 	}
// 	if err := database.DB.Create(&user).Error; err != nil {
// 		return c.JSON(500, err)
// 	}
// 	return c.JSON(200, user)
// }
