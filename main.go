package main // import "shiryaev"

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"image/png"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	f, err := os.Open("shiryaev.png")
	if err != nil {
		panic(err)
	}

	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}

	var phase1 []byte
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			// Read color, and convert it from 48-bit (?) to 24-bit
			r, g, b, _ := img.At(x, y).RGBA()
			r = r >> 8
			g = g >> 8
			b = b >> 8

			phase1 = append(phase1, byte(r), byte(g), byte(b))
		}
	}
	fmt.Println("Phase 1: ", phase1)

	// Two rounds of Deflate decompressing
	phase2, err := ioutil.ReadAll(flate.NewReader(flate.NewReader(bytes.NewReader(phase1))))
	if err != nil {
		panic(err)
	}
	fmt.Println("Phase 2: ", phase2)

	// Optional, to support Russian symbols
	phase3, err := charmap.Windows1251.NewDecoder().Bytes(phase2)
	if err != nil {
		panic(err)
	}
	fmt.Println("Phase 3: ", string(phase3))

	// Decoding Base64 string
	phase4 := make([]byte, base64.StdEncoding.DecodedLen(len(phase3)))
	_, err = base64.StdEncoding.Decode(phase4, phase3)
	if err != nil {
		panic(err)
	}
	fmt.Println("Phase 4: ", string(phase4))

	// Trim and decode binary string representation to bytes
	var phase5 []byte
	for _, b := range strings.Split(string(phase4), " ") {
		b = strings.Trim(b, string([]byte{0}))
		i, err := strconv.ParseInt(b, 2, 64)
		if err != nil {
			panic(err)
		}
		phase5 = append(phase5, byte(i))
	}
	fmt.Println("Phase 5: ", string(phase5))

	// Finally convert bytes aka string to int
	phase6, err := strconv.Atoi(string(phase5))
	if err != nil {
		panic(err)
	}

	// Print UTC time from UNIX
	fmt.Println("Phase 6: ", time.Unix(int64(phase6), 0).UTC())
}
