package cmd

// MIT License

// Copyright (c) 2019 7onetella

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var attrToNames = map[string]string{
	"uid":       "User ID",
	"sn":        "Last Name",
	"homePhone": "Home Phone",
	"cn":        "Full Name",
}

var userShowCmd = &cobra.Command{
	Use:     "show <uid>",
	Short:   "Shows user",
	Example: "show smithj",
	Run: func(cmd *cobra.Command, args []string) {

		lc := readLdapConfig()

		uid := args[0]

		attributes, groups, _ := findLdapUser(lc, uid)
		if attributes == nil {
			Failure("user does not exists")
			return
		}

		//Log(fmt.Sprintf("%+v", attributes))
		Print("\n")
		for k, v := range attributes {
			name, ok := attrToNames[k]
			if !ok {
				name = k
			}
			Log(fmt.Sprintf("%-20s: %s\n", name, v))
		}
		Log(fmt.Sprintf("%-20s: %s", "Groups", strings.Join(groups, ",")))
	},
}

func init() {
	userCmd.AddCommand(userShowCmd)
}
