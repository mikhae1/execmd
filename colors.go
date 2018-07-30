package execmd

import (
	"hash/fnv"

	"github.com/fatih/color"
)

var defaultColors = []color.Attribute{
	color.FgCyan,
	color.FgYellow,
	color.FgMagenta,
	color.FgBlue,

	color.FgHiGreen,
	color.FgHiYellow,
	color.FgHiBlue,
	color.FgHiMagenta,
	color.FgHiCyan,
}

func Color(str string) (coloredStr string) {
	hash := fnv.New32a()

	hash.Write([]byte(str))

	colorAtrr := defaultColors[hash.Sum32()%uint32(len(defaultColors))]

	coloredStr = color.New(colorAtrr).SprintFunc()(str)

	return
}

func ColorErr(str string) string {
	return color.New(color.FgRed).SprintFunc()(str)
}

func ColorOK(str string) string {
	return color.New(color.FgGreen).SprintFunc()(str)
}

func ColorStrong(str string) string {
	return color.New(color.Bold).SprintFunc()(str)
}
