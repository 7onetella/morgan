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

	"github.com/7onetella/morgan/internal/cryptow/ssha"
	"github.com/spf13/cobra"
)

var userCreateCmdFirst string
var userCreateCmdLast string
var userCreateCmdHomePhone string

const ldifCreateTemplate = `version: 1

dn: {{.Dn}}
objectclass: inetOrgPerson
objectclass: organizationalPerson
objectclass: person
objectclass: top
cn: {{.Cn}}
sn: {{.Sn}}
uid: {{.UID}}
homePhone: {{.HomePhone}}
userPassword: {{.UserPassword}}`

var userCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Creates user",
	Example: "--first John --last Smith --home_phone 777-555-1234",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := readLdapConfig()
		ldap := NewLDAP(cfg)

		uid := strings.ToLower(userCreateCmdLast) + strings.ToLower(userCreateCmdFirst)[0:1]
		passwordHash, err := ssha.Encode([]byte(cfg.Password))
		ExitOnError(err, "hashing password")

		// clean up slack home phone formatting. e.g. <tel:777-555-1234|777-555-1234>
		if strings.Contains(userCreateCmdHomePhone, "|") {
			p := strings.Split(userCreateCmdHomePhone, "|")[1]
			p = strings.Replace(p, ">", "", 1)
			userCreateCmdHomePhone = p
		}

		ldapEntry := LdapEntry{
			Dn:           fmt.Sprintf(cfg.Dn, uid),
			Cn:           userCreateCmdFirst + " " + userCreateCmdLast,
			Sn:           userCreateCmdLast,
			UID:          uid,
			HomePhone:    userCreateCmdHomePhone,
			UserPassword: string(passwordHash),
		}

		attributes, _, _ := ldap.findLdapUser(uid)
		if attributes != nil {
			Info("attributes: " + fmt.Sprintf("%+v", attributes))
			Failure("user already exists")
			return
		}

		path, err := createLDIF(applyTemplate(ldapEntry, ldifCreateTemplate))
		ExitOnError(err, "creating temporary ldif file")

		ldap.do("ldapadd", "creating user", path)

		cleanUpLDIF(path)

		Success("creating user")
	},
}

func init() {
	userCmd.AddCommand(userCreateCmd)

	flags := userCreateCmd.Flags()

	flags.StringVar(&userCreateCmdFirst, "first", "", "first name")
	flags.StringVar(&userCreateCmdLast, "last", "", "last name")
	flags.StringVar(&userCreateCmdHomePhone, "home_phone", "", "home phone")

}
