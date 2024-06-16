package main

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go"
)

// Replace with your secret key for signing JWTs
var jwtSecret = []byte("your_secret_key")

func main() {
    r := gin.Default()
    r.POST("/forgot-password", forgotPassword)
    r.GET("/reset-password", showResetPasswordForm)
    r.POST("/reset-password", resetPassword)

    r.Run(":8080")
}

func forgotPassword(c *gin.Context) {
    var request struct {
        Email string `json:"email"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // Generate JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email": request.Email,
        "exp":   time.Now().Add(time.Hour).Unix(), // Token expires in 1 hour
    })

    // Sign the token with the secret key
    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    // Send email with link containing the token (e.g., http://yourdomain.com/reset-password?token={token})
    sendEmail(request.Email, tokenString)

    c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}

func sendEmail(email, token string) {
    // Implement your email sending logic here
    // This is a placeholder function
    println("Sending email to:", email, "with token:", token)
}


func showResetPasswordForm(c *gin.Context) {
    token := c.Query("token")

    // Parse and validate the token
    claims := jwt.MapClaims{}
    _, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
        return
    }

    // Optionally, you can check if the token has expired here
    // Note: Expiry check can also be done automatically by the JWT library

    // Serve the HTML form for password reset
    // In a real application, you would render a template here
    c.String(http.StatusOK, "Reset your password (token: %s)", token)
}


func resetPassword(c *gin.Context) {
    var request struct {
        Token    string `json:"token"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // Parse and validate the token
    claims := jwt.MapClaims{}
    _, err := jwt.ParseWithClaims(request.Token, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
        return
    }

    // Optionally, you can check if the token has expired here
    // Note: Expiry check can also be done automatically by the JWT library

    email := claims["email"].(string)

    // Update the user's password in the database
    // This is a placeholder function
    updatePassword(email, request.Password)

    c.JSON(http.StatusOK, gin.H{"message": "Password has been reset"})
}

func updatePassword(email, password string) {
    // Implement your password update logic here
    // This is a placeholder function
    println("Updating password for:", email, "to:", password)
}
