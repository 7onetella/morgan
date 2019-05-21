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

	"github.com/spf13/cobra"
)

const ldifDelteTemplate = `{{.Dn}} `

var userDeleteCmd = &cobra.Command{
	Use:     "delete <uid>",
	Short:   "Deletes user",
	Example: "delete smithj",
	Run: func(cmd *cobra.Command, args []string) {

		lc := readLdapConfig()

		uid := args[0]

		ldapEntry := LdapEntry{
			Dn: fmt.Sprintf(lc.Dn, uid),
		}

		attributes, _, _ := findLdapUser(lc, uid)
		if attributes == nil {
			Info("attributes: " + fmt.Sprintf("%+v", attributes))
			Failure("user does not exists")
			return
		}

		path, err := createLDIF(applyTemplate(ldapEntry, ldifDelteTemplate))
		ExitOnError(err, "creating temporary ldif file")

		cmds := []string{"ldapdelete", "-h", lc.Host, "-p", lc.Port, "-D", lc.AdminDn, "-w", lc.AdminPwd, "-f", path}
		// Debug(strings.Join(cmds, " "))

		executeLdapCmd("deleting user", cmds)

		cleanUpLDIF(path)

		Success("deleting user")
	},
}

func init() {
	userCmd.AddCommand(userDeleteCmd)

	// flags := userDeleteCmd.Flags()
}
