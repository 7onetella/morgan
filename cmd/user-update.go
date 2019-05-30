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

var userUpdateCmdDebug bool
var userUpdateCmdHomePhone string

const ldifUpdateTemplate = `dn: {{.Dn}}
changetype: modify
replace: homePhone
homePhone: {{.HomePhone}}`

var userUpdateCmd = &cobra.Command{
	Use:     "update <uid>",
	Short:   "Updates user",
	Example: "update smithj",
	Run: func(cmd *cobra.Command, args []string) {

		lc := readLdapConfig()
		ldap := NewLDAP(lc)

		uid := args[0]

		ldapEntry := LdapEntry{
			Dn:        fmt.Sprintf(lc.Dn, uid),
			HomePhone: userUpdateCmdHomePhone,
		}

		attributes, _, _ := ldap.findLdapUser(uid)
		if attributes == nil {
			Info("attributes: " + fmt.Sprintf("%+v", attributes))
			Failure("user does not exists")
			return
		}

		path, err := createLDIF(applyTemplate(ldapEntry, ldifUpdateTemplate))
		ExitOnError(err, "creating temporary ldif file")

		cmds := []string{"ldapmodify", "-h", lc.Host, "-p", lc.Port, "-D", lc.AdminDn, "-w", lc.AdminPwd, "-f", path}
		argsmasked := []string{"ldapmodify", "-h", lc.Host, "-p", lc.Port, "-D", "*****", "-w", "*****", "-f", path}
		Debug(strings.Join(argsmasked, " "))

		executeLdapCmd("updating user", cmds)

		cleanUpLDIF(path)

		Success("updating user")
	},
}

func init() {
	userCmd.AddCommand(userUpdateCmd)

	flags := userUpdateCmd.Flags()

	flags.StringVar(&userUpdateCmdHomePhone, "home_phone", "", "home phone")
}
