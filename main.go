package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

/* This program generates a random 128x128 image using random.org HTTP API */

func main() {
	const width, height = 128, 128
	const size = width * height * 3

	rnd, err := getRandomNumbers(size, 0, 255)
	if err != nil {
		log.Fatal(err)
	}
	if len(rnd) < size {
		log.Fatal("Not enough random numbers generated!")
	}

	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b := rnd[i], rnd[i+1], rnd[i+2]
			img.Set(x, y, color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255})
			i += 3
		}
	}

	f, err := os.Create("image.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func getRandomNumbers(count int, min, max int) ([]int, error) {
	// Make sure each request is at most for 10000 numbers
	// Break it down into multiple requests if necessary
	var result []int
	for len(result) < count {
		err := makeRandomRequest(minimum(10000, count-len(result)), min, max, &result)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

func minimum(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func makeRandomRequest(count int, min, max int, appendTo *[]int) error {
	// Real server: "https://www.random.org/integers/"
	// Mock server: "http://localhost:8080/integers/"
	server := "https://www.random.org/integers/"
	url := fmt.Sprintf("%s?num=%d&min=%d&max=%d&col=1&base=10&format=plain&rnd=new",
		server, count, min, max)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Status %v %s", resp.Status, body)
	}
	err = nil
	for err == nil {
		var x, n int
		n, err = fmt.Fscanf(resp.Body, "%d\n", &x)
		if n == 1 {
			*appendTo = append(*appendTo, x)
		}
	}
	return nil
}
