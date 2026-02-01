package middleware

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	APP_KEY = "www.topgoer.com"
)

// TokenHandler是我们获取用户名和密码的处理程序，如果有效，则返回用于将来请求的令牌。
func TokenHandler(userId int) (string, error) {
	log.Printf("开始为用户：%d 生成token", userId)
	// 颁发一个有限期一小时的证书
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		"iat":    time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(APP_KEY))
	log.Printf("用户：%d 生成token：%v", userId, token)
	return tokenString, err
}
func ParseToken(tokenString string) (jwt.MapClaims, string) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(APP_KEY), nil
	})
	if token.Valid {
		fmt.Println("You look nice today")
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
		return claims, ""
	} else {
		fmt.Println(err)
	}
	return nil, ""
}
func RefreshTokenHandler(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":    time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(APP_KEY))
	return tokenString, err
}
func ETagStaticMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentVersion := "v1.0.0"
		etag := fmt.Sprintf(`W/"%x"`, md5.Sum([]byte(contentVersion)))
		clientETag := c.GetHeader("If-None-Match")
		if clientETag == etag {
			c.AbortWithStatus(http.StatusNotModified)
			return
		}
		c.Header("ETag", etag)
		c.Header("Cache-Control", "no-cache")
		c.Next()
	}
}
