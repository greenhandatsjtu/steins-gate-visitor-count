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
)

type Count struct {
	Value int `json:"value"`
}

func cacheImages(digits *[]image.Image) {
	for i := range *digits {
		file, _ := os.Open(fmt.Sprintf("digits/%d.png", i))
		defer file.Close()
		(*digits)[i], _, _ = image.Decode(file)
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
	req, _ := http.NewRequest("GET", "https://api.countapi.xyz/hit/visitor-badge/"+key, nil)
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
	for i, _ := range count {
		index, _ := strconv.Atoi(count[i : i+1])
		draw.Draw(img, image.Rect(i*200, 0, 200*length, 200), digits[index], digits[index].Bounds().Min, draw.Over)
	}
	return img
}
