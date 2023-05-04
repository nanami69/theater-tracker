package geocoding

import (
	"net/http"
	"fmt"
    "encoding/json"
    "net/url"
)

func GetLatLng(address string) (string, string, error) {
    // APIに送信するリクエストを作成
    apiURL := "https://msearch.gsi.go.jp/address-search/AddressSearch?q=" + url.QueryEscape(address)
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