package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func registerCinema(c echo.Context) error {
	// データベースの登録と緯度経度算出処理
	return c.JSON(http.StatusOK, map[string]string{
		"message": "test",
	})
}

func main() {
	e := echo.New()

	// CORSの設定(デプロイ前に設定を見直す)
	e.Use(middleware.CORS())
	// ルーティングの設定
	e.GET("/users/:id", getUser)
	e.POST("/register-cinema", registerCinema)
	// サーバーの起動
	e.Logger.Fatal(e.Start(":8080"))
}
