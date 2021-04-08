package output

var reset = "\033[0m"
var danger = "\033[31m"
var success = "\033[32m"
var info = "\033[36m"
var warning = "\033[0;33m"

func Success(text string) {
	outputColor(success, text)
}

func Danger(text string) {
	outputColor(danger, text)
}

func Info(text string) {
	outputColor(info, text)
}

func Warning(text string) {
	outputColor(warning, text)
}

func outputColor(color, text string) {
	println(color + text + reset)
}
