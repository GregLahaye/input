package input

import (
	"bufio"
	"fmt"
	"github.com/GregLahaye/yogurt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
)

func Rune() rune {
	// set terminal to raw mode so we can read one character at a time
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return ' '
	}
	defer terminal.Restore(int(os.Stdin.Fd()), oldState)

	reader := bufio.NewReader(os.Stdin)

	r, _, err := reader.ReadRune()
	if err != nil {
		return ' '
	}

	return r
}

func String() string {
	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(s)
}

func Multiline(prompt string) string {
	fmt.Println(prompt)

	s := ""
	reader := bufio.NewReader(os.Stdin)
	for {
		if i, err := reader.ReadString('\n'); err != nil {
			return s
		} else if strings.TrimSpace(i) == "" {
			return s
		} else {
			s += i
		}
	}
}

func Confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s ", prompt)
		s, err := reader.ReadString('\n')
		if err != nil {
			return false
		}
		c := s[0]

		if c == 'y' || c == 'Y' {
			return true
		} else if c == 'n' || c == 'N' {
			return false
		}
	}
}

func Select(options []string) int {
	length := len(options)

	// use a string builder so everything is displayed at once
	var s strings.Builder

	// disable cursor so it's not blinking
	s.WriteString(yogurt.DisableCursor)

	// for each option, add a selection box and display it
	for _, o := range options {
		s.WriteString(" [ ] ")
		s.WriteString(o)
		s.WriteString("\n")
	}

	// move cursor to column three so it is in selection box
	s.WriteString(yogurt.SetColumn(3))
	// move cursor back to first option
	s.WriteString(yogurt.CursorUp(length))
	// display initial 'x' in first option's box
	s.WriteString("x")
	// move cursor backwards one so it's back in selection box
	s.WriteString(yogurt.CursorBackward(1))

	// display string
	fmt.Print(s.String())

	// reset for cleanup string
	s.Reset()

	index := 0 // index of currently selected option
	done := false
	for !done {
		c := Rune()
		switch c {
		case 3: // CTRL-C
			done = true
		case 13: // Enter
			done = true
		case 'j':
			if index+1 < length {
				// clear 'x' in current selection box
				fmt.Print(" ")
				// move cursor back into selection box
				fmt.Print(yogurt.CursorBackward(1))
				// move cursor down into next selection box
				fmt.Print(yogurt.CursorDown(1))
				// display 'x' in current selection box
				fmt.Print("x")
				// move cursor backwards one so it's back in selection box
				fmt.Print(yogurt.CursorBackward(1))

				index++
			}
		case 'k':
			if index > 0 {
				// clear 'x' in current selection box
				fmt.Print(" ")
				// move cursor back into selection box
				fmt.Print(yogurt.CursorBackward(1))
				// move cursor up into previous selection box
				fmt.Print(yogurt.CursorUp(1))
				// display 'x' in current selection box
				fmt.Print("x")
				// move cursor backwards one so it's back in selection box
				fmt.Print(yogurt.CursorBackward(1))

				index--
			}
		}
	}

	// set to first column
	s.WriteString(yogurt.SetColumn(0))
	if index > 0 {
		// move up to first index if we have moved down
		s.WriteString(yogurt.CursorUp(index))
	}

	// for each option, clear line
	for j := 0; j < length; j++ {
		s.WriteString(yogurt.ClearLine)
		s.WriteString(yogurt.CursorDown(1))
	}

	// move cursor up to top
	s.WriteString(yogurt.CursorUp(length))

	// display string
	fmt.Print(s.String())

	// re-enable cursor
	fmt.Print(yogurt.EnableCursor)

	return index
}
