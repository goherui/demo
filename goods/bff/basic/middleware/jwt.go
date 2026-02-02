package middleware

import (
	"crypto/md5"
	"errors"
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
func NegotiationCacheMiddleware(str1 string) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentVersion := "v0.0.1"
		etag := fmt.Sprintf(`W/"%x"`, md5.Sum([]byte(contentVersion)))
		lastModified := time.Now().Add(-1 * time.Hour).Format(http.TimeFormat)
		clientETag := c.GetHeader("If-None-Match")
		clientModified := c.GetHeader("If-Modified-Since")
		if clientETag == etag || clientModified == lastModified {
			c.AbortWithStatus(http.StatusNotModified)
			return
		}
		c.Header(str1, etag)
		c.Header("Last-Modified", lastModified)
		c.Header("Cache-Control", "max-age=3600")
		c.Next()
	}
}
func CreateToken(tokenString string) (string, error) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(APP_KEY), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("token格式错误")
	}
	userId, ok := claims["userId"].(int)
	if ok {
		return "", errors.New("无法获取用户id")
	}

	return TokenHandler(userId)
}
