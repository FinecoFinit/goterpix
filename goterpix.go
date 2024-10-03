package main

import (
	"fmt"
	"golang.org/x/term"
	"image/gif"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

type Pixel struct {
	R, G, B, A int
}

var clearscr map[string]func()

func main() {
	var (
		backVT, foreVT string
		back, fore     Pixel
	)
	// Args
	path := os.Args[1]
	delayArg := os.Args[2]
	delay, _ := strconv.Atoi(delayArg)
	// Commands for terminal clean
	clearscr = make(map[string]func()) //Initialize it
	clearscr["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}
	clearscr["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}
	// Input file
	inputFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func(inputFile *os.File) {
		err := inputFile.Close()
		if err != nil {
			panic(err)
		}
	}(inputFile)
	g, err := gif.DecodeAll(inputFile)
	if err != nil {
		panic(err)
	}
	// Declaring types and sizes
	HeightStart, WidthStart := 0, 0
	Width, Height := g.Config.Width, g.Config.Height
	frames := make([]string, len(g.Image))
	terX, terY, err := term.GetSize(0)
	if err != nil {
		panic(err)
	}
	// Set size
	if terY*2 < Height {
		HeightStart = (Height - terY*2) / 2
		Height = Height - ((Height - terY*2) / 2)
	}
	if terX < Width {
		WidthStart = (Width - terX) / 2
		Width = Width - ((Width - terX) / 2)
	}
	// Magic
	for f := 0; f < len(g.Image); f++ {
		for y := HeightStart; y < Height; y += 2 {
			for x := WidthStart; x < Width; x += 1 {
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
	for i := 0; i < len(frames); i++ {
		time.Sleep(time.Duration(delay) * time.Millisecond)
		CallClear()
		print(frames[i])
	}
}
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}
func CallClear() {
	value, ok := clearscr[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
