package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func registerCinema(c echo.Context) error {
	// データベースの登録と緯度経度算出処理
	var req struct {
		Name     string `json:"name"`
		Address  string `json:"address"`
		PhotoURL string `json:"photo_url"`
		Comment  string `json:"comment"`
	}
	if err := c.Bind(&req); err != nil {
		return err
	}

	lat, lng, err := getLatLng(req.Address)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"name": req.Name,
		"lat": lat,
		"lng": lng,
		"comment": req.Comment,
	})
}

func getLatLng(address string) (string, string, error) {
    // APIに送信するリクエストを作成
    apiURL := "https://msearch.gsi.go.jp/address-search/AddressSearch?q=" + url.QueryEscape(address)
    // values := make(url.Values)
    // values.Set("q", address)
    req, err := http.NewRequest("GET", apiURL, nil)
    if err != nil {
        return "", "", err
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    // リクエストを送信
    client := new(http.Client)
    resp, err := client.Do(req)
    if err != nil {
        return "", "", err
    }
    defer resp.Body.Close()

    // レスポンスをパース
    var res []struct {
        Geometry struct {
            Coordinates []float64 `json:"coordinates"`
        } `json:"geometry"`
    }
    err = json.NewDecoder(resp.Body).Decode(&res)
    if err != nil {
        return "", "", err
    }

    // レスポンスから緯度経度を取得
    if len(res) == 0 {
        return "", "", fmt.Errorf("no result")
    }
    longitude := fmt.Sprintf("%f", res[0].Geometry.Coordinates[0])
    latitude := fmt.Sprintf("%f", res[0].Geometry.Coordinates[1])

    return latitude, longitude, nil
}

func main() {
	e := echo.New()

	// CORSの設定(デプロイ前に設定を見直す)
	e.Use(middleware.CORS())
	// ルーティングの設定
	e.POST("/register-cinema", registerCinema)
	// サーバーの起動
	e.Logger.Fatal(e.Start(":8080"))
}
