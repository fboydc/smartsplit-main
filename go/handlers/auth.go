package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fboydc/smartsplit-main/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your_secret_key")

// AuthHandlers contains all authentication-related HTTP handlers
type AuthHandlers struct {
	db *sql.DB
}

// NewAuthHandlers creates a new AuthHandlers instance
func NewAuthHandlers(db *sql.DB) *AuthHandlers {
	return &AuthHandlers{
		db: db,
	}
}

// LoginHandler handles user login requests
func (h *AuthHandlers) LoginHandler(c *gin.Context) {
	uname := c.PostForm("user")
	passwd := c.PostForm("password")

	auth, userid, plaidToken, err := h.authenticateUser(uname, passwd)
	if err != nil || !auth {
		if err == sql.ErrNoRows || !auth {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid credentials",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error: Could not authenticate user",
			})
		}
		return
	}

	token, err := h.generateJWT(userid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error: Could not generate token",
		})
		return
	}

	response := models.AuthResponse{
		Message:    "Login successful",
		Token:      token,
		PlaidToken: plaidToken,
		UserID:     userid,
		Username:   uname,
	}

	c.JSON(http.StatusOK, response)
}

// AuthMiddleware validates JWT tokens for protected routes
func (h *AuthHandlers) AuthMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		// Note: You might want to handle access token here too
		accessToken := c.GetHeader("AccessToken")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		valid, err := h.validateJWT(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to validate token"})
			c.Abort()
			return
		}
		if token == "" || !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
			c.Abort()
			return
		}

		c.Set("AccessToken", accessToken)
		c.Next()
	})
}

// AuthHandler handles auth endpoint requests
func (h *AuthHandlers) AuthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "authenticated",
	})
}

// Private helper methods

func (h *AuthHandlers) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (h *AuthHandlers) comparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePlainPwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlainPwd)
	return err == nil
}

func (h *AuthHandlers) authenticateUser(username string, password string) (bool, string, string, error) {
	var hashedPassword, plaidAccessToken, userid sql.NullString
	err := h.db.QueryRow(`SELECT user_id, password_hash, plaid_access_token FROM "Users" WHERE username = $1`, username).Scan(&userid, &hashedPassword, &plaidAccessToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, "", "", err
		}
		return false, "", "", err
	}

	return h.comparePasswords(hashedPassword.String, password), userid.String, plaidAccessToken.String, nil
}

func (h *AuthHandlers) generateJWT(userid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": userid,
		"exp":    time.Now().Add(time.Minute * 15).Unix(),
	})

	return token.SignedString(jwtSecret)
}

func (h *AuthHandlers) validateJWT(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return false, err
	}

	return token.Valid, nil
}

// SaveAccessToken saves the access token to the database
func (h *AuthHandlers) SaveAccessToken(accessToken string, user string) (bool, error) {
	_, err := h.db.Exec(`UPDATE "Users" SET plaid_access_token = $1 where username = $2`, accessToken, user)
	if err != nil {
		return false, err
	}
	log.Printf("Access token updated for user: %s", user)
	return true, nil
}
