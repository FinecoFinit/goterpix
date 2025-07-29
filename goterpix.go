package main

import (
	"atomicgo.dev/cursor"
	"fmt"
	"image"
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

type Tab struct {
	Up Pixel
	Dn Pixel
}

type TabsRow struct {
	Tabs []Tab
}

type Frame struct {
	TabsRows []TabsRow
}

type Frames struct {
	Frames []string
}

var ClearScr map[string]func()

func main() {
	path := os.Args[1]
	delayArg := os.Args[2]
	delay, _ := strconv.Atoi(delayArg)

	ClearScr = make(map[string]func()) //Initialize it
	ClearScr["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}
	ClearScr["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}

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

	//HeightStart, WidthStart := 0, 0
	//Width, Height := g.Config.Width, g.Config.Height
	//fd := int(os.Stdout.Fd())
	//terX, terY, err := terminal.GetSize(fd)
	//if err != nil {
	//	panic(err)
	//}

	//if terY*2 < Height {
	//	HeightStart = (Height - terY*2) / 2
	//	Height = Height - ((Height - terY*2) / 2)
	//}
	//if terX < Width {
	//	WidthStart = (Width - terX) / 2
	//	Width = Width - ((Width - terX) / 2)
	//}

	var frames Frames

	for _, f := range g.Image {
		frames.Frames = append(frames.Frames, BuildAnsi(BuildRows(f)))
	}
	cursor.Hide()
	for _, frame := range frames.Frames {
		CallClear()
		_, err = os.Stdout.Write([]byte(frame))
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	cursor.Show()
}

func BuildRows(f *image.Paletted) Frame {
	var frame Frame
	for pr := 0; pr < f.Rect.Dy(); pr += 2 {
		var TabsRow TabsRow
		for p := 0; p < f.Rect.Dx(); p++ {
			var Tab Tab
			Tab.Up = RgbaToPixel(f.At(p, pr).RGBA())
			Tab.Dn = RgbaToPixel(f.At(p, pr+1).RGBA())
			TabsRow.Tabs = append(TabsRow.Tabs, Tab)
		}
		frame.TabsRows = append(frame.TabsRows, TabsRow)
	}
	return frame
}

func BuildAnsi(f Frame) string {
	var fr string
	for _, tr := range f.TabsRows {
		var tempRow string
		for _, tab := range tr.Tabs {
			switch {
			case tab.Up.A == 0 && tab.Dn.A == 0:
				tempRow += " "
			case tab.Up.A > 0 && tab.Dn.A == 0:
				UP := fmt.Sprint("\033[38;2;", tab.Up.R, ";", tab.Up.G, ";", tab.Up.B, "m")
				tempRow += fmt.Sprint(UP, "▀\033[0m")
			case tab.Up.A == 0 && tab.Dn.A > 0:
				DN := fmt.Sprint("\033[38;2;", tab.Dn.R, ";", tab.Dn.G, ";", tab.Dn.B, "m")
				tempRow += fmt.Sprint(DN, "▄\033[0m")
			case tab.Up.A > 0 && tab.Dn.A > 0:
				UP := fmt.Sprint("\033[48;2;", tab.Up.R, ";", tab.Up.G, ";", tab.Up.B, "m")
				DN := fmt.Sprint("\033[38;2;", tab.Dn.R, ";", tab.Dn.G, ";", tab.Dn.B, "m")
				tempRow += fmt.Sprint(UP, DN, "▄\033[0m")
			}
		}
		tempRow += "\n"
		fr += tempRow
	}
	return fr
}

func RgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func CallClear() {
	value, ok := ClearScr[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
