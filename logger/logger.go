package logger

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
)

var (
	// Green is Success Logger.
	Green *log.Logger

	// Blue is Note Logger.
	Blue *log.Logger

	// Red is Error Logger.
	Red *log.Logger
)

const (
	red = uint8(iota + 91)
	green
	yellow
	blue
	magenta
)

func color(color uint8, str string) string {
	return fmt.Sprintf("\x1b[1;%dm%s\x1b[0m", color, str)
}

func init() {
	Blue = log.New(os.Stdout, color(blue, "\t[*] "), 0)
	Green = log.New(os.Stdout, color(green, "[+] "), 0)
	Red = log.New(os.Stdout, color(red, "[!] "), 0)

	logo, _ := base64.StdEncoding.DecodeString("ICAgX19fX18gICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgLl9fICAgICAgICAg\nX19fX19fX18gICAgICAgIAogIC8gIF8gIFwgICBfX19fX18gX19fX19fX19fXyAgICBfX19fX18g\nX19fX198X198IF9fX18gIC8gIF9fX19fLyAgX19fXyAgCiAvICAvX1wgIFwgLyAgX19fLy8gIF9f\nX3xfXyAgXCAgLyAgX19fLy8gIF9fXy8gIHwvICAgIFwvICAgXCAgX19fIC8gIF8gXCAKLyAgICB8\nICAgIFxfX18gXCBcX19fIFwgLyBfXyBcX1xfX18gXCBcX19fIFx8ICB8ICAgfCAgXCAgICBcX1wg\nICggIDxfPiApClxfX19ffF9fICAvX19fXyAgPl9fX18gID5fX19fICAvX19fXyAgPl9fX18gID5f\nX3xfX198ICAvXF9fX19fXyAgL1xfX19fLyAKICAgICAgICBcLyAgICAgXC8gICAgIFwvICAgICBc\nLyAgICAgXC8gICAgIFwvICAgICAgICBcLyAgICAgICAgXC8gICAgICAgIAo=")
	fmt.Println(color(red, string(logo)))
}
