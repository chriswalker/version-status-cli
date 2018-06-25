package output

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

type Version struct {
	ServiceName    string
	StagingVersion string
	ProdVersion    string
}

type Outputter interface {
	Output([]string, []Version)
}

type StdOutputter struct{}

func NewStdOutputter() *StdOutputter {
	return &StdOutputter{}
}

const (
	Reset  string = "\x1b[0000m"
	Blue   string = "\x1b[0034m"
	Red    string = "\x1b[0031m"
	Yellow string = "\x1b[0033m"
	Green  string = "\x1b[0032m"
)

func (s *StdOutputter) Output(contexts []string, versions []Version) {
	fmt.Printf("Version differences between %s and %s\n\n", contexts[0], contexts[1])

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Little hacky; tabwriter doesn't recognise colour escape sequences and so doesn't calculate
	// tabbing correctly; ensure the header row has the same number of escape sequences as the
	// data rows, by wrapping a couple of columns in the 'reset' code
	fmt.Fprintf(w, "\tService\t%s\t%s\t\n", colour(Reset, contexts[0]), colour(Reset, contexts[1]))

	for _, ver := range versions {
		c := getColour(ver)
		fmt.Fprintf(w, "\t%s\t%s\t%s\t\n", ver.ServiceName, colour(c, ver.StagingVersion), colour(c, ver.ProdVersion))
	}

	w.Flush()

	fmt.Println()
}

func getColour(version Version) string {
	lv, err := stringSplitToIntSplit(strings.Split(strings.Split(version.StagingVersion, "-")[0], "."))
	if err != nil {
		return Blue
	}
	rv, err := stringSplitToIntSplit(strings.Split(strings.Split(version.ProdVersion, "-")[0], "."))
	if err != nil {
		return Blue
	}

	switch {
	case lv[0] != rv[0]:
		return Red
	case lv[1] == rv[1] && lv[2] == rv[2]:
		return Green
	case lv[1]-rv[1] > 1 || lv[1]-rv[1] < -1:
		return Red
	case lv[1]-rv[1] < 2 || lv[1]-rv[1] > -2:
		return Yellow
	case lv[2] != rv[2]:
		return Yellow
	}

	return Blue
}

func stringSplitToIntSplit(strings []string) (ints []int, err error) {
	ints = make([]int, len(strings))
	for i, s := range strings {
		ints[i], err = strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
	}

	return ints, err
}

// colour wraps the supplied string with the supplied colour and reset codes
func colour(colour, s string) string {
	return fmt.Sprintf("%s%s%s", colour, s, Reset)
}
