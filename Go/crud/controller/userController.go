package controller

import (
	"github.com/cdfmlr/crud/model"
	"github.com/cdfmlr/crud/orm"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func SignUp(c *gin.Context) {
	var body struct {
		Email    string `json:"Email"` // Ensure fields are exported by starting with a capital letter
		Password string `json:"Password"`
	}

	// Use 'err := c.Bind(&body); err != nil' to correctly handle the binding error
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return // Return here to stop function execution after an error
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return // Properly handling the error and stopping the execution
	}

	user := model.User{Email: body.Email, Password: string(hash)}
	result := orm.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return // Properly handling the error and stopping the execution
	}

	// This will be the only response sent to the client if everything above succeeds
	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}
