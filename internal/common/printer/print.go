package printer

import (
	"fmt"
	"io"
	"os"
	"workshop/internal/common/printer/styles"

	"github.com/charmbracelet/x/term"
)

// See: https://github.com/hay-kot/scaffold/blob/main/internal/printer/printer.go
// See: https://github.com/charmbracelet/lipgloss

// Console provides simple, styled console output
// using lipgloss for consistent branding.
type Console struct {
	writer io.Writer
}

// New returns a new Console instance that writes to
// the given io.Writer.
func New(w io.Writer) *Console {
	return &Console{writer: w}
}

// SetWriter changes the output writer for the console.
// Useful for redirecting output in tests or file logs.
func (c *Console) SetWriter(w io.Writer) {
	c.writer = w
}

// Print writes plain, unstyled text to the console.
func (c *Console) Print(msg string) {
	c.write(msg)
}

// Println writes plain, unstyled text to the console, with a linebreak.
func (c *Console) Println(msg string) {
	c.write(msg)
	c.LineBreak()
}

// Printf writes plain, unstyled text to the console,  with formatting.
func (c *Console) Printf(msg string, args ...any) {
	c.write(fmt.Sprintf(msg, args...))
}

// Success prints a success message with a checkmark icon and success color.
func (c *Console) Success(msg string) {
	c.Println(styles.Success.Render(fmt.Sprintf("%s %s", styles.IconSuccess, msg)))
}

// Successf formats and prints a success message.
func (c *Console) Successf(format string, args ...any) {
	c.Success(fmt.Sprintf(format, args...))
}

// Error prints an error message with a cross icon and error color.
func (c *Console) Error(msg string) {
	c.Println(styles.Error.Render(fmt.Sprintf("%s %s", styles.IconError, msg)))
}

// Errorf formats and prints an error message.
func (c *Console) Errorf(format string, args ...any) {
	c.Error(fmt.Sprintf(format, args...))
}

// Info prints an informational message with an info icon and color.
func (c *Console) Info(msg string) {
	c.Println(styles.Info.Render(fmt.Sprintf("%s %s", styles.IconInfo, msg)))
}

// Infof formats and prints an informational message.
func (c *Console) Infof(format string, args ...any) {
	c.Info(fmt.Sprintf(format, args...))
}

// Warn prints a warning message with a warning icon and color.
func (c *Console) Warn(msg string) {
	c.Println(styles.Warn.Render(fmt.Sprintf("%s %s", styles.IconWarn, msg)))
}

// Warnf formats and prints a warning message.
func (c *Console) Warnf(format string, args ...any) {
	c.Warn(fmt.Sprintf(format, args...))
}

// LineBreak prints \n to the writer.
func (c *Console) LineBreak() {
	c.write("\n")
}

// isTerminal reports whether the console writes to an interactive
// terminal, as opposed to a pipe, file or buffer.
func (c *Console) isTerminal() bool {
	f, ok := c.writer.(*os.File)
	return ok && term.IsTerminal(f.Fd())
}

func (c *Console) write(s string) {
	if c.writer == nil { // Guard check
		return
	}
	_, _ = io.WriteString(c.writer, s)
}
