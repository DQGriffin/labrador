package console

import (
	"fmt"
	"os"

	"github.com/DQGriffin/labrador/internal/cli/styles"
)

var debugPrefix = "[DEBUG] "
var warnPrefix = "[WARN] "
var errorPrefix = "[ERROR] "

var isColorEnabled = false
var isDebugOutputEnabled = false

func SetColorEnabled(value bool) {
	isColorEnabled = value
}

func SetDebugOutputEnabled(value bool) {
	isDebugOutputEnabled = value
}

func Debug(args ...interface{}) {
	if !isDebugOutputEnabled {
		return
	}

	text := fmt.Sprint(args...)
	debugText := debugPrefix + text
	if isColorEnabled {
		fmt.Fprintln(os.Stderr, styles.Primary.Render(debugText))
	} else {
		fmt.Fprintln(os.Stderr, debugText)
	}
}

func Debugf(format string, args ...interface{}) {
	if !isDebugOutputEnabled {
		return
	}
	output := fmt.Sprintf(format, args...)
	formatted := debugPrefix + output
	if isColorEnabled {
		fmt.Fprintln(os.Stderr, styles.Primary.Render(formatted))
	} else {
		fmt.Fprintln(os.Stderr, formatted)
	}
}

func Info(args ...interface{}) {
	text := fmt.Sprint(args...)
	if isColorEnabled {
		fmt.Println(styles.Primary.Render(text))
	} else {
		fmt.Println(text)
	}
}

func Infof(format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	if isColorEnabled {
		fmt.Println(styles.Primary.Render(output))
	} else {
		fmt.Println(output)
	}
}

func Warn(args ...interface{}) {
	text := fmt.Sprint(args...)
	warnText := warnPrefix + text
	if isColorEnabled {
		fmt.Fprintln(os.Stderr, styles.Warn.Render(warnText))
	} else {
		fmt.Fprint(os.Stderr, warnText)
	}
}

func Warnf(format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	formatted := warnPrefix + output
	if isColorEnabled {
		fmt.Fprintln(os.Stderr, styles.Warn.Render(formatted))
	} else {
		fmt.Fprintln(os.Stderr, formatted)
	}
}

func Error(args ...interface{}) {
	text := fmt.Sprint(args...)
	errorText := errorPrefix + text
	if isColorEnabled {
		fmt.Fprintln(os.Stderr, styles.Error.Render(errorText))
	} else {
		fmt.Fprintln(os.Stderr, errorText)
	}
}

func Errorf(format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	formatted := errorPrefix + output
	if isColorEnabled {
		fmt.Fprintln(os.Stderr, styles.Error.Render(formatted))
	} else {
		fmt.Fprintln(os.Stderr, formatted)
	}
}

func Fatal(text string) {
	Error(text)
	os.Exit(1)
}

func Fatalf(format string, args ...interface{}) {
	Errorf(format, args...)
	os.Exit(1)
}
