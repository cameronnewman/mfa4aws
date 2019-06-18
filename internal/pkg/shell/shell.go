package shell

import (
	"fmt"
	"io"
)

//PrintVars prints vars array to io.writer
func PrintVars(out io.Writer, vars []string) {
	for _, x := range vars {
		fmt.Fprintf(out, "%s\n", x)
	}
}
