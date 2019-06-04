// Copyright © 2019 Christian Weichel
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/32leaves/packer/pkg/twod"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

var printSVG bool

// twodCmd represents the 2D command
var twodCmd = &cobra.Command{
	Use:   "2d <problem.json>",
	Short: "Solves a 2D packing problem",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fn := args[0]
		fc, err := ioutil.ReadFile(fn)
		if err != nil {
			fmt.Printf("cannot read problem: %v\n", err)
			os.Exit(1)
		}

		var problem twod.Problem
		err = json.Unmarshal(fc, &problem)
		if err != nil {
			fmt.Printf("cannot read problem: %v\n", err)
			os.Exit(1)
		}

		sol, err := twod.NextFitDecreasingHeight(problem)
		if err != nil {
			fmt.Printf("cannot solve problem: %v\n", err)
			os.Exit(1)
		}

		if printSVG {
			fmt.Println(sol.SVG())
			return
		}

		out, err := json.MarshalIndent(sol, "", "  ")
		if err != nil {
			fmt.Printf("cannot marshal solution: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(out))
	},
}

func init() {
	rootCmd.AddCommand(twodCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// 2dCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	twodCmd.Flags().BoolVar(&printSVG, "svg", false, "Print solution as SVG")
}
