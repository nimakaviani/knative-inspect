package logs

import (
	"fmt"
	"log"
	"os"
	"strings"

	. "github.com/cppforlife/go-cli-ui/ui/table"
)

type DumpUI struct {
	filter string
}

func NewDumpUI(filter string) *DumpUI {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	return &DumpUI{filter: filter}
}

func (o *DumpUI) Close() {
	log.SetFlags(log.Flags() & (log.Ldate | log.Ltime))
	log.SetOutput(os.Stderr)
}

func (o *DumpUI) ErrorLinef(pattern string, args ...interface{}) {
	log.Printf(pattern, args...)
}

func (o *DumpUI) PrintLinef(pattern string, args ...interface{}) {
	log.Printf(pattern, args...)
}

func (o *DumpUI) BeginLinef(pattern string, args ...interface{}) {
	log.Printf(pattern, args...)
}

func (o *DumpUI) EndLinef(pattern string, args ...interface{}) {
	log.Printf(pattern, args...)
}

func (o *DumpUI) PrintBlock(data []byte) { // takes []byte to avoid string copy
	logLine := string(data)
	if o.filter == "" {
		log.Print(logLine)
		return
	}

	if strings.Contains(logLine, fmt.Sprintf(`"level":"%s"`, o.filter)) {
		log.Print(logLine)
	}
}

func (o *DumpUI) PrintErrorBlock(string) {
}

func (o *DumpUI) PrintTable(Table) {
}

func (o *DumpUI) AskForText(label string) (string, error) {
	return "", nil
}

func (o *DumpUI) AskForChoice(label string, options []string) (int, error) {
	return -1, nil
}

func (o *DumpUI) AskForPassword(label string) (string, error) {
	return "", nil
}

// AskForConfirmation returns error if user doesnt want to continue
func (o *DumpUI) AskForConfirmation() error { return nil }
func (o *DumpUI) IsInteractive() bool       { return false }
func (o *DumpUI) Flush()                    {}
