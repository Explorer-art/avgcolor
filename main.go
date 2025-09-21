package main

import (
	"fmt"
	"os"
	"io"
	"path/filepath"
	"strings"
	"image"
	"image/png"
	"image/jpeg"
)

type Pixel struct {
	R int
	G int
	B int
	A int
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func getAverageColor(file io.Reader) (Pixel, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return Pixel{}, err
	}

	var avg_r, avg_g, avg_b, avg_a int

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	if width == 0 || height == 0 {
		return Pixel{}, fmt.Errorf("invalid image size")
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := rgbaToPixel(img.At(x, y).RGBA())
			avg_r += pixel.R
			avg_g += pixel.G
			avg_b += pixel.B
			avg_a += pixel.A
		}
	}

	return Pixel{avg_r / (height * width), avg_g / (height * width), avg_b / (height * width), avg_a / (height * width)}, nil
}

func parseFlags() map[string]string {
	flags := make(map[string]string, 2)
	var mode string
	
	for i := 1; i < len(os.Args); i++ {
		if mode == "" && (os.Args[i] == "-f" || os.Args[i] == "--format") {
			mode = "-f"
		} else if mode == "-f" {
			flags["-f"] = os.Args[i]
			mode = ""
		}
	}

	return flags
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <image> [flags]\n", os.Args[0]);
		os.Exit(0)
	}

	path := os.Args[1]

	flags := parseFlags()

	if _, ok := flags["-f"]; !ok {
		ext := strings.ToLower(filepath.Ext(path))
			
		if ext == "" {
			fmt.Fprintln(os.Stderr, "error: unknown image format")
			os.Exit(1)
		}
		
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			fmt.Fprintf(os.Stderr, "error: unknown %s format\n", ext[1:])
			os.Exit(1)
		}
		
		if ext == ".png" {
			image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
		} else if ext == ".jpg" || ext == ".jpeg" {
			image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
		}
	} else {
		if strings.ToLower(flags["-f"]) == "png" {
			image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
		} else if strings.ToLower(flags["-f"]) == "jpg" || strings.ToLower(flags["-f"]) == "jpeg" {
			image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
		}
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	defer file.Close()

	pixel, err := getAverageColor(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: image could not be decoded")
		os.Exit(1)
	}

	r, g, b := pixel.R, pixel.G, pixel.B

	fmt.Printf("\n \x1b[48;2;%d;%d;%dm        \x1b[0m\t\tAverage color\n \x1b[48;2;%d;%d;%dm        \x1b[0m\t\tRGB: %d, %d, %d\n \x1b[48;2;%d;%d;%dm        \x1b[0m\t\tHEX: #%02x%02x%02x\n\n", r, g, b, r, g, b, r, g, b, r, g, b, r, g, b)
}
