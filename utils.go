package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/nfnt/resize"
)

type Count struct {
	Value int `json:"value"`
}

func cacheImages(digits *[]image.Image) {
	cacheOneImage := func(no int) {
		file, _ := os.Open(fmt.Sprintf("digits/%d.png", no))
		defer file.Close()
		(*digits)[no], _, _ = image.Decode(file)
	}
	for i := range *digits {
		cacheOneImage(i) // to avoid resource leak
	}
}

func generateMd5(id string) (string, error) {
	w := md5.New()
	if _, err := io.WriteString(w, id); err != nil {
		return "", err
	}
	res := fmt.Sprintf("%x", w.Sum(nil))
	return res, nil
}

func updateCounter(key string) string {
	req, _ := http.NewRequest("GET", "https://api.countapi.xyz/hit/steins-gate-visitor-count/"+key, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	var count Count
	if err = json.Unmarshal(body, &count); err != nil {
		log.Println(err)
		return ""
	}
	return strconv.Itoa(count.Value)
}

func generateImage(digits []image.Image, count string) image.Image {
	length := len(count)
	img := image.NewNRGBA(image.Rect(0, 0, 200*length, 200))
	for i := range count {
		index, _ := strconv.Atoi(count[i : i+1])
		draw.Draw(img, image.Rect(i*200, 0, 200*length, 200), digits[index], digits[index].Bounds().Min, draw.Over)
	}
	return img
}

// resizeImage resize image to specified ratio
func resizeImage(img image.Image, ratio float64) image.Image {
	width := uint(float64(img.Bounds().Max.X-img.Bounds().Min.X) * ratio)
	height := uint(float64(img.Bounds().Max.Y-img.Bounds().Min.Y) * ratio)
	return resize.Resize(width, height, img, resize.Lanczos3)
}
