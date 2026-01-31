package terminal

import (
	"fmt"
	"time"
)

const (
	COLORBLACK    int = 0
	COLORBLUE     int = 4
	COLORWHITE    int = 15
	COLORDARKGRAY int = 235
	COLORYELLOW   int = 11
)

// MoveToColumn returns an ANSI escape to move the cursor to a column.
func MoveToColumn(col int) string {
	return fmt.Sprintf("\033[%dG", col)
}

// ResetColor returns the ANSI escape to reset colors.
func ResetColor() string {
	return "\033[39m\033[49m"
}

// SetBgColor returns the ANSI escape to set background color.
func SetBgColor(bgcolor int) string {
	return fmt.Sprintf("\033[48;5;%dm", bgcolor)
}

// SetFgColor returns the ANSI escape to set foreground color.
func SetFgColor(bgcolor int) string {
	return fmt.Sprintf("\033[38;5;%dm", bgcolor)
}

// GetStrikeTroughText wraps text with strikethrough ANSI codes.
func GetStrikeTroughText(text string) string {
	return fmt.Sprintf("\033[9m%s\033[m", text)
}

// GetHeadlineText formats a headline with background color.
func GetHeadlineText(text string) string {
	return fmt.Sprintf("%s%s%s", SetBgColor(COLORBLUE), text, ResetColor())
}

// GetLabelText formats a label with background and foreground colors.
func GetLabelText(text string) string {
	return fmt.Sprintf("%s%s%-12s%s", SetBgColor(COLORDARKGRAY), SetFgColor(COLORYELLOW), text, ResetColor())
}

// GetTimeString formats a Unix timestamp as a date string.
func GetTimeString(ts int64) string {
	t := time.Unix(ts, 0)
	return t.Format("02.01.2006 15:04")
}
