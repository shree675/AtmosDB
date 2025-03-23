package util

import (
	"fmt"

	"github.com/fatih/color"
)

func PrintError(val string) {
	fmt.Println(color.RedString("[ERROR]"), val)
}

func PrintRed(val string) {
	color.Red(val)
}

func PrintGreen(val string) {
	color.Green(val)
}

func PrintYellow(val string) {
	color.Yellow(val)
}

func PrintBlue(val string) {
	color.Blue(val)
}

func GetMagentaStr(val string) string {
	return color.MagentaString(val)
}

func GetYellowStr(val string) string {
	return color.YellowString(val)
}

func PrintGray(val string) {
	color.RGB(150, 150, 150).Println(val)
}
