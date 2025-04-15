package version

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Embedded by --ldflags on build time
// Versioning should follow the SemVer guidelines
// https://semver.org/
var (
	// Version is the main version at the moment.
	branch    string
	commit    string // the git commit that the binary was built on
	author    string
	date      string
	goVersion string
)

var cmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("branch     : ", branch)
		fmt.Println("author     : ", author)
		fmt.Println("date       : ", date)
		fmt.Println("git commit : ", commit)
		fmt.Println("go version : ", goVersion)
	},
}

func NewCommand() *cobra.Command {
	return cmd
}
