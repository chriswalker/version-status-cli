package output

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
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

func (s *StdOutputter) Output(contexts []string, versions []Version) {
	w := tabwriter.NewWriter(os.Stdout, 0, 20, 0, '\t', 0)
	fmt.Fprintf(w, "Service\t%s\t%s\t\n", contexts[0], contexts[1])
	for _, ver := range versions {
		fn := getColourFunc(ver)
		fmt.Fprintf(w, "%s\t%s", ver.ServiceName, fn("%s\t%s\t\n", ver.StagingVersion, ver.ProdVersion))
	}
	w.Flush()
}

func getColourFunc(version Version) func(format string, a ...interface{}) string {
	lv, err := stringSplitToIntSplit(strings.Split(strings.Split(version.StagingVersion, "-")[0], "."))
	if err != nil {
		return color.New(color.FgBlue).SprintfFunc()
	}
	rv, err := stringSplitToIntSplit(strings.Split(strings.Split(version.ProdVersion, "-")[0], "."))
	if err != nil {
		return color.New(color.FgBlue).SprintfFunc()
	}

	switch {
	case lv[0] != rv[0]:
		return color.New(color.FgRed).SprintfFunc()
	case lv[1] == rv[1] && lv[2] == rv[2]:
		return color.New(color.FgGreen).SprintfFunc()
	case lv[1]-rv[1] > 1 || lv[1]-rv[1] < -1:
		return color.New(color.FgRed).SprintfFunc()
	case lv[1]-rv[1] < 2 || lv[1]-rv[1] > -2:
		return color.New(color.FgYellow).SprintfFunc()
	case lv[2] != rv[2]:
		return color.New(color.FgYellow).SprintfFunc()
	}

	return color.New(color.FgBlue).SprintfFunc()
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
