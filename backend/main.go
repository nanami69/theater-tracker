package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func getUser(c echo.Context) error {
	// id := c.Param("id")
	// ダミーデータとしてユーザーを返す
	user := &User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}
	return c.JSON(http.StatusOK, user)
}

func main() {
	e := echo.New()
	// ルーティングの設定
	e.GET("/users/:id", getUser)
	// サーバーの起動
	e.Logger.Fatal(e.Start(":8080"))
}
