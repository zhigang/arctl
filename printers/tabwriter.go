package printers

import (
	"fmt"
	"io"

	"text/tabwriter"

	"github.com/olekukonko/tablewriter"
)

// IOStreams provides the standard names for iostreams.  This is useful for embedding and for unit testing.
// Inconsistent and different names make it hard to read and review code
type IOStreams struct {
	// In think, os.Stdin
	In io.Reader
	// Out think, os.Stdout
	Out io.Writer
	// ErrOut think, os.Stderr
	ErrOut io.Writer
}

const (
	tabwriterMinWidth = 6
	tabwriterWidth    = 4
	tabwriterPadding  = 3
	tabwriterPadChar  = ' '
	tabwriterFlags    = tabwriter.DiscardEmptyColumns
)

// GetNewTableWriter returns a tablewriter that translates tabbed columns in input into properly aligned text.
func GetNewTableWriter(output io.Writer) *tablewriter.Table {
	t := tablewriter.NewWriter(output)
	t.SetAutoWrapText(false)
	t.SetAutoFormatHeaders(true)
	t.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	t.SetFooterAlignment(tablewriter.ALIGN_LEFT)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetCenterSeparator("")
	t.SetColumnSeparator("")
	t.SetRowSeparator("")
	t.SetHeaderLine(false)
	t.SetBorder(false)
	t.SetTablePadding("\t") // pad with tabs
	// t.SetNoWhiteSpace(true)
	return t
}

// Print prints information with default TableWriter.
func Print(output io.Writer, headers []string, values ...[]string) {
	t := GetNewTableWriter(output)
	t.SetHeader(headers)
	t.AppendBulk(values)
	t.Render()
	fmt.Println()
}

// PrintTable prints information with a new TableWriter.
func PrintTable(output io.Writer, headers []string, values ...[]string) {
	t := tablewriter.NewWriter(output)
	t.SetHeader(headers)
	t.AppendBulk(values)
	t.Render()
}

// GetNewTabWriter returns a tabwriter that translates tabbed columns in input into properly aligned text.
func GetNewTabWriter(output io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(output, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, tabwriterFlags)
}

// PrintRows prints rows with default tabwriter.
func PrintRows(out io.Writer, rows []string) {
	for i := 0; i < len(rows); i++ {
		PrintValue(out, rows[i])
	}
	fmt.Fprintln(out)
}

// PrintValue prints value with default tabwriter.
func PrintValue(out io.Writer, value interface{}) {
	fmt.Fprintf(out, "%v\t", value)
}

// PrintObj prints name and value with default tabwriter.
func PrintObj(out io.Writer, name string, value interface{}) {
	fmt.Fprintf(out, "%s:\t%v\t\n", name, value)
}
