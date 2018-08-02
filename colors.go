package execmd

import (
	"hash/fnv"

	fcolor "github.com/fatih/color"
)

var defaultColors = []fcolor.Attribute{
	fcolor.FgCyan,
	fcolor.FgYellow,
	fcolor.FgMagenta,
	fcolor.FgBlue,

	fcolor.FgHiGreen,
	fcolor.FgHiYellow,
	fcolor.FgHiBlue,
	fcolor.FgHiMagenta,
	fcolor.FgHiCyan,
}

func color(str string) (coloredStr string) {
	hash := fnv.New32a()

	hash.Write([]byte(str))

	colorAtrr := defaultColors[hash.Sum32()%uint32(len(defaultColors))]

	coloredStr = fcolor.New(colorAtrr).SprintFunc()(str)

	return
}

func colorErr(str string) string {
	return fcolor.New(fcolor.FgRed).SprintFunc()(str)
}

func colorOK(str string) string {
	return fcolor.New(fcolor.FgGreen).SprintFunc()(str)
}

func colorStrong(str string) string {
	return fcolor.New(fcolor.Bold).SprintFunc()(str)
}
