package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/base64"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	"github.com/nanami69/theater-tracker/backend/pkg/geocoding"
)

type Config struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	DBName   string `mapstructure:"dbname"`
}

func saveToDB(name string, address string, latitude string, longitude string, photo []byte) error {
    // 設定ファイルのパスを指定する
	viper.SetConfigFile("./config/config.yml")

	// 設定ファイルを読み込む
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %s", err))
	}

	// 設定ファイルの内容を構造体にマッピングする
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal config file: %s", err))
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.User, config.Password, config.Host, config.Port, config.DBName)

    // DBに接続
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return err
    }
    defer db.Close()

    // INSERT文を作成して実行
    query := "INSERT INTO theaters (name, address, latitude, longitude, photo) VALUES (?, ?, ?, ?, ?)"
    _, err = db.Exec(query, name, address, latitude, longitude, photo)
    if err != nil {
        return err
    }

    return nil
}

func registerCinema(c echo.Context) error {
	// データベースの登録と緯度経度算出処理
    var req struct {
        Name     string `json:"name"`
        Address  string `json:"address"`
        Comment  string `json:"comment"`
        FileData []byte `json:"photo"`
    }
	var fileBase64 = ""
	var blobData = []byte("")

    // ファイルを受け取る
    file, err := c.FormFile("photo")
    if err != nil {
        if err == http.ErrMissingFile {
            // ファイルが送信されていない場合の処理
        } else {
            // その他のエラーが発生した場合の処理
            return fmt.Errorf("resieve：%s", err.Error())
        }
    } else {
        // ファイルが正常に受け取れた場合の処理
        f, err := file.Open()
        if err != nil {
            return fmt.Errorf("not opne：%s", err.Error())
        }
        defer f.Close()

        data, err := ioutil.ReadAll(f)
        if err != nil {
            return fmt.Errorf("read file：%s", err.Error())
        }
        // Base64エンコード
    	fileBase64 = base64.StdEncoding.EncodeToString(data)
		blobData = []byte(fileBase64)
        fmt.Println("files: uploaded")
    }

    // その他のリクエストデータを受け取る
    if err := c.Bind(&req); err != nil {
		fmt.Println("Bind")
        return fmt.Errorf("bind：%s", err.Error())
    }

	lat, lng, err := geocoding.GetLatLng(c.FormValue("address"))
	if err != nil {
		return err
	}

	saveToDB(c.FormValue("name"), c.FormValue("address"), lat, lng, blobData)

	responseMap := map[string]string{
		"name": c.FormValue("name"),
		"lat": lat,
		"lng": lng,
		"comment": c.FormValue("comment"),
	}
	
	if fileBase64 != "" {
		responseMap["photo"] = fileBase64
	}
	
	return c.JSON(http.StatusOK, responseMap)
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
