package console

import (
	"fmt"
	"os"

	"github.com/mgutz/ansi"
)

func Success(msg string) {
	colorOut(msg, "green")
}

func Error(msg string) {
	colorOut(msg, "red")
}

func Warning(msg string) {
	colorOut(msg, "yellow")
}

// Exit 打印一条报错消息，并退出 os.Exit(1)
func Exit(msg string) {
	Error(msg)
	os.Exit(1)
}

func ExitIf(err error) {
	if err != nil {
		Exit(err.Error())
	}
}

func colorOut(message string, color string) {
	fmt.Fprintln(os.Stdout, ansi.Color(message, color))
}
