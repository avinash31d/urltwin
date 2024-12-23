package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/avinash31d/urltwin/config"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

func GoogleLogin(c *fiber.Ctx) error {
	url := config.GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(url)
}

func GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(http.StatusBadRequest).SendString("Code not found")
	}

	token, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Println("Error during token exchange:", err)
		return c.Status(http.StatusInternalServerError).SendString("Failed to exchange token")
	}

	client := config.GoogleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Println("Error getting user info:", err)
		return c.Status(http.StatusInternalServerError).SendString("Failed to get user info")
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Println("Error decoding user info:", err)
		return c.Status(http.StatusInternalServerError).SendString("Failed to parse user info")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token.AccessToken,
		HTTPOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		Expires:  time.Now().Add(time.Until(token.Expiry)),
	})

	return c.JSON(fiber.Map{
		"message":  "Login successful",
		"userInfo": userInfo,
	})
}

func Logout(c *fiber.Ctx) error {
	c.ClearCookie("auth_token")
	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})
}
