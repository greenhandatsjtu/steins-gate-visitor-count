package main

import (
	"bytes"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	port := os.Getenv("port")
	if port == "" {
		port = "8080"
	}

	digits := make([]image.Image, 10)
	cacheImages(&digits)

	e := echo.New()
	e.Use(middleware.Recover())
	//e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")

		m, err := generateMd5(id)
		if err != nil {
			log.Println(err)
			return c.NoContent(http.StatusBadRequest)
		}

		count := updateCounter(m)
		if count == "" {
			log.Println("Fetch visitor count error.")
			return c.NoContent(http.StatusInternalServerError)
		}

		img := generateImage(digits, count)
		if v := c.QueryParam("ratio"); len(v) != 0 {
			if ratio, err := strconv.ParseFloat(v, 64); err == nil {
				img = resizeImage(img, ratio)
			} else {
				print(err)
			}
		}
		buf := new(bytes.Buffer)
		err = png.Encode(buf, img)
		if err != nil {
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}

		expireTime := time.Now().Add(-10 * time.Minute).String()
		c.Response().Header().Add("Expires", expireTime)
		c.Response().Header().Add("Cache-Control", "no-cache,max-age=0,no-store,s-maxage=0,proxy-revalidate")

		return c.Blob(http.StatusOK, "image/png", buf.Bytes())
	})

	e.Logger.Fatal(e.Start(":" + port))
}
