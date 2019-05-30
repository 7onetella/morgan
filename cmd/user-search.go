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
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var userSearchCmd = &cobra.Command{
	Use:     "search <Partial Full Name>",
	Short:   "Searches users",
	Example: `search "John Smith"`,
	Run: func(cmd *cobra.Command, args []string) {

		cfg := readLdapConfig()
		ldap := LDAP{cfg}

		filter := args[0]

		users, _ := ldap.findUsers("*" + filter + "*")
		if users == nil || len(users) == 0 {
			Println("  \u2022 search result empty")
			return
		}

		//Log(fmt.Sprintf("%+v", attributes))
		Print("\n")

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Home#", "Email"})

		for _, user := range users {
			table.Append([]string{user["cn"], user["homePhone"], user["mail"]})
		}
		table.Render() // Send output
	},
}

func init() {
	userCmd.AddCommand(userSearchCmd)
}
