/*
Copyright © 2023 J.Kushibiki
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	var patterns []string

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
			// check if patterns are specified
			if len(patterns) == 0 {
				return fmt.Errorf("no patterns specified")
			}

			// with files
			if len(args) > 0 {
				// open files
				var files []io.Reader
				for _, arg := range args {
					f, err := os.Open(arg)
					if err != nil {
						return err
					}
					defer f.Close()
					files = append(files, f)
				}
				// concat files
				reader := io.MultiReader(files...)
				// run gipp
				return Run(reader, os.Stdout, os.Stderr, patterns)
			}

			// without files
			return Run(os.Stdin, os.Stdout, os.Stderr, patterns)
		},
	}

	cmd.Flags().StringSliceVarP(&patterns, "pattern", "e", []string{}, "pattern")

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)

	return cmd
}

func Run(in io.Reader, out, eout io.Writer, ps []string) error {
	// load patterns
	patterns := make([]Pattern, len(ps))
	for i, p := range ps {
		pattern, err := ParsePattern(p)
		if err != nil {
			return err
		}
		patterns[i] = pattern
	}

	// read input stream line by line
	sc := bufio.NewScanner(in)
	for sc.Scan() {
		line := sc.Text()
		// parse line
		ip, err := ParseIp(line)
		if err != nil {
			continue
		}

		// match patterns
		for _, pattern := range patterns {
			if pattern.Match(ip) {
				fmt.Fprintln(out, line)
			}
		}
	}

	return nil
}

func Execute() {
	err := NewRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
