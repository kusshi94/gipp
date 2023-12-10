/*
Copyright © 2023 J.Kushibiki
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gipp [flags] [-e pattern] [file ...]",
		Short: "IP Prefix/Suffix Version of grep",
		Long: `The gipp utility searches any given IP address list files, selecting lines that match one or more patterns.
The pattern is written in an extended cidr notation that allows suffixes to be expressed.

following are examples of the pattern:
	192.168.100.0/24
	0.0.0.1/-8
	::abcd:01ff:fe00:0/-64/24`,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(os.Stdin, os.Stdout, os.Stderr, args)
		},
	}

	var pattern string
	cmd.Flags().StringVarP(&pattern, "pattern", "e", "", "pattern")

	return cmd
}

func run(in io.Reader, out, eout io.Writer, files []string) error {
	fmt.Fprintln(out, "Hello, gipp!")
	return nil
}

func Execute() {
	err := NewRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
