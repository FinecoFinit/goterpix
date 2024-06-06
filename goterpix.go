package main

import (
	"fmt"
	"image/gif"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

var clear map[string]func() // Map for storing clear funcs
var backVT, foreVT string   // Strings for pixel building
var back, fore Pixel        // Pixel color

func main() {
	// Args
	path := os.Args[1]
	delayArg := os.Args[2]
	delay, _ := strconv.Atoi(delayArg)

	// Commands for terminal clean
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	// Input file
	inputFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()
	g, err := gif.DecodeAll(inputFile)
	if err != nil {
		panic(err)
	}

	// Declaring types
	Width, Height := g.Config.Width, g.Config.Height
	frames := make([]string, len(g.Image))

	// Magic
	for f := 0; f < len(g.Image); f++ {
		for y := 0; y < Height; y += 2 {
			for x := 0; x < Width; x += 1 {

				// Symbol structure in ANSI escape code, back is upper half and fore is bottom half
				back = rgbaToPixel(g.Image[f].RGBA64At(x, y).RGBA())
				fore = rgbaToPixel(g.Image[f].RGBA64At(x, y+1).RGBA())
				backVT = fmt.Sprint("\033[48;2;", back.R, ";", back.G, ";", back.B, "m")
				foreVT = fmt.Sprint("\033[38;2;", fore.R, ";", fore.G, ";", fore.B, "m")

				if fore.A == 0 && back.A == 0 {
					frames[f] += " "
				}
				if fore.A == 0 && back.A > 0 {
					backVT = fmt.Sprint("\033", "[38;2;", back.R, ";", back.G, ";", back.B, "m")
					frames[f] += fmt.Sprint(backVT, "▀\033[0m")
				}
				if fore.A > 0 && back.A == 0 {
					frames[f] += fmt.Sprint(foreVT, "▄\033[0m")
				}
				if fore.A > 0 && back.A > 0 {
					frames[f] += fmt.Sprint(backVT, foreVT, "▄\033[0m")
				}
			}
			frames[f] += "\n"
		}
	}
	// Draw magic
	for i := 0; i < len(frames); i++ {
		time.Sleep(time.Duration(delay) * time.Millisecond)
		CallClear()
		print(frames[i])
	}
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

// Pixel struct
type Pixel struct {
	R, G, B, A int
}

// Clearing function
func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
