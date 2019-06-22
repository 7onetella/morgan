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
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/7onetella/morgan/internal/execw"
	"github.com/alecthomas/template"
	"github.com/google/uuid"
	// ldap "github.com/jtblin/go-ldap-client"
	"crypto/tls"
	"errors"
	"fmt"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"

	"gopkg.in/ldap.v2"
)

var userCmd = &cobra.Command{
	Use:   "user [command]",
	Short: "Automation for users",
	Long:  `Automation for users`,
}

func init() {
	rootCmd.AddCommand(userCmd)
}

// LdapConfig ldap config
type LdapConfig struct {
	Base       string `yaml:"base"`
	Group      string `yaml:"group"`
	User       string `yaml:"user"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Dn         string `yaml:"dn"`
	AdminDn    string `yaml:"binddn"`
	AdminPwd   string `yaml:"bindpwd"`
	Password   string `yaml:"password"`
	Attributes string `yaml:"attributes"`
}

// LdapEntry ldap entry
type LdapEntry struct {
	Dn           string
	Cn           string
	Sn           string
	UID          string
	UserPassword string
	HomePhone    string
}

func createLDIF(b []byte) (string, error) {

	Info(string(b))

	path := "/tmp/" + uuid.New().String()
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return path, err
	}

	return path, nil
}

func readLdapConfig() LdapConfig {
	home, err := homedir.Dir()
	ExitOnError(err, "finding home")

	lc := LdapConfig{}
	data, err := ioutil.ReadFile(home + "/.morgan.ldap.yml")
	ExitOnError(err, "reading ldap yml")
	yaml.Unmarshal(data, &lc)

	return lc
}

func applyTemplate(ldapEntry LdapEntry, tpl string) []byte {
	tmpl, err := template.New("tmp").Parse(tpl)
	ExitOnError(err, "parsing template")

	var b bytes.Buffer
	err = tmpl.Execute(&b, ldapEntry)
	ExitOnError(err, "processing template")

	return b.Bytes()
}

func cleanUpLDIF(path string) {
	err := os.Remove(path)
	ExitOnError(err, "removing temporary ldif file")
}

// deprecated. use ldapCreate or ldapModify instead
func executeLdapCmd(msg string, cmds []string) {
	stdout, errout, err := execw.Exec(cmds)
	ExitOnErrorWithDetail(err, msg, stdout, errout)
}

// LDAP handles various ldap operations
type LDAP struct {
	cfg LdapConfig
}

// NewLDAP constructor
func NewLDAP(cfg LdapConfig) LDAP {
	return LDAP{cfg}
}

func (l LDAP) do(cmd, msg string, path string) {
	cfg := l.cfg
	cmds := []string{cmd, "-h", cfg.Host, "-p", l.cfg.Port, "-D", cfg.AdminDn, "-w", cfg.AdminPwd, "-f", path}
	argsmasked := []string{cmd, "-h", cfg.Host, "-p", cfg.Port, "-D", "*****", "-w", "*****", "-f", path}
	Info(strings.Join(argsmasked, " "))

	stdout, errout, err := execw.Exec(cmds)
	ExitOnErrorWithDetail(err, msg, stdout, errout)
}

func (l LDAP) findLdapUser(uid string) (map[string]string, []string, error) {
	cfg := l.cfg
	port, _ := strconv.Atoi(cfg.Port)

	client := &LDAPClient{
		Base:         cfg.Base,
		Host:         cfg.Host,
		Port:         port,
		UseSSL:       false,
		BindDN:       cfg.AdminDn,
		BindPassword: cfg.AdminPwd,
		UserFilter:   cfg.User,
		GroupFilter:  cfg.Group,
		Attributes:   strings.Split(cfg.Attributes, ","),
	}
	// It is the responsibility of the caller to close the connection
	defer client.Close()

	_, user, err := client.SearchBy(uid)
	if err != nil {
		return nil, nil, err
	}

	groups, err := client.GetGroupsOfUser(uid)
	if err != nil {
		return nil, nil, err
	}

	return user, groups, nil
}

func (l LDAP) findUsers(filter string) ([]map[string]string, error) {
	cfg := l.cfg
	port, _ := strconv.Atoi(cfg.Port)

	client := &LDAPClient{
		Base:         cfg.Base,
		Host:         cfg.Host,
		Port:         port,
		UseSSL:       false,
		BindDN:       cfg.AdminDn,
		BindPassword: cfg.AdminPwd,
		UserFilter:   cfg.User,
		GroupFilter:  cfg.Group,
		Attributes:   strings.Split(cfg.Attributes, ","),
	}
	// It is the responsibility of the caller to close the connection
	defer client.Close()

	users, err := client.SearchByFullName(filter)
	if err != nil {
		return nil, err
	}

	return users, nil
}

/*
	All credit for the following code goes to https://github.com/jtblin/go-ldap-client
*/

// LDAPClient ldap client
type LDAPClient struct {
	Attributes         []string
	Base               string
	BindDN             string
	BindPassword       string
	GroupFilter        string // e.g. "(memberUid=%s)"
	Host               string
	ServerName         string
	UserFilter         string // e.g. "(uid=%s)"
	Conn               *ldap.Conn
	Port               int
	InsecureSkipVerify bool
	UseSSL             bool
	SkipTLS            bool
	ClientCertificates []tls.Certificate // Adding client certificates
}

// Connect connects to the ldap backend.
func (lc *LDAPClient) Connect() error {
	if lc.Conn == nil {
		var l *ldap.Conn
		var err error
		address := fmt.Sprintf("%s:%d", lc.Host, lc.Port)
		if !lc.UseSSL {
			l, err = ldap.Dial("tcp", address)
			if err != nil {
				return err
			}

			// Reconnect with TLS
			if !lc.SkipTLS {
				err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
				if err != nil {
					return err
				}
			}
		} else {
			config := &tls.Config{
				InsecureSkipVerify: lc.InsecureSkipVerify,
				ServerName:         lc.ServerName,
			}
			if lc.ClientCertificates != nil && len(lc.ClientCertificates) > 0 {
				config.Certificates = lc.ClientCertificates
			}
			l, err = ldap.DialTLS("tcp", address, config)
			if err != nil {
				return err
			}
		}

		lc.Conn = l
	}
	return nil
}

// Close closes the ldap backend connection.
func (lc *LDAPClient) Close() {
	if lc.Conn != nil {
		lc.Conn.Close()
		lc.Conn = nil
	}
}

// SearchBy searches the user against the ldap backend.
func (lc *LDAPClient) SearchBy(uid string) (bool, map[string]string, error) {
	err := lc.Connect()
	if err != nil {
		return false, nil, err
	}

	// First bind with a read only user
	if lc.BindDN != "" && lc.BindPassword != "" {
		err := lc.Conn.Bind(lc.BindDN, lc.BindPassword)
		if err != nil {
			return false, nil, err
		}
	}

	attributes := append(lc.Attributes, "dn")
	// Search for the given cn
	searchRequest := ldap.NewSearchRequest(
		lc.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(lc.UserFilter, uid),
		attributes,
		nil,
	)

	sr, err := lc.Conn.Search(searchRequest)
	if err != nil {
		return false, nil, err
	}

	if len(sr.Entries) < 1 {
		return false, nil, errors.New("User does not exist")
	}

	if len(sr.Entries) > 1 {
		return false, nil, errors.New("Too many entries returned")
	}

	user := map[string]string{}
	for _, attr := range lc.Attributes {
		user[attr] = sr.Entries[0].GetAttributeValue(attr)
	}

	return true, user, nil
}

// SearchByFullName searches the user against the ldap backend by full name
func (lc *LDAPClient) SearchByFullName(filter string) ([]map[string]string, error) {
	err := lc.Connect()
	if err != nil {
		return nil, err
	}

	// First bind with a read only user
	if lc.BindDN != "" && lc.BindPassword != "" {
		err := lc.Conn.Bind(lc.BindDN, lc.BindPassword)
		if err != nil {
			return nil, err
		}
	}

	attributes := append(lc.Attributes, "dn")
	// Search for the given cn
	searchRequest := ldap.NewSearchRequest(
		lc.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(cn=%s)", filter),
		attributes,
		nil,
	)

	sr, err := lc.Conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	if len(sr.Entries) < 1 {
		return nil, errors.New("User does not exist")
	}

	users := []map[string]string{}
	for _, entry := range sr.Entries {
		user := map[string]string{}
		for _, attr := range lc.Attributes {
			user[attr] = entry.GetAttributeValue(attr)
		}
		users = append(users, user)
	}

	return users, nil
}

// Authenticate authenticates the user against the ldap backend.
func (lc *LDAPClient) Authenticate(username, password string) (bool, map[string]string, error) {
	err := lc.Connect()
	if err != nil {
		return false, nil, err
	}

	// First bind with a read only user
	if lc.BindDN != "" && lc.BindPassword != "" {
		err := lc.Conn.Bind(lc.BindDN, lc.BindPassword)
		if err != nil {
			return false, nil, err
		}
	}

	attributes := append(lc.Attributes, "dn")
	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		lc.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(lc.UserFilter, username),
		attributes,
		nil,
	)

	sr, err := lc.Conn.Search(searchRequest)
	if err != nil {
		return false, nil, err
	}

	if len(sr.Entries) < 1 {
		return false, nil, errors.New("User does not exist")
	}

	if len(sr.Entries) > 1 {
		return false, nil, errors.New("Too many entries returned")
	}

	userDN := sr.Entries[0].DN
	user := map[string]string{}
	for _, attr := range lc.Attributes {
		user[attr] = sr.Entries[0].GetAttributeValue(attr)
	}

	// Bind as the user to verify their password
	err = lc.Conn.Bind(userDN, password)
	if err != nil {
		return false, user, err
	}

	// Rebind as the read only user for any further queries
	if lc.BindDN != "" && lc.BindPassword != "" {
		err = lc.Conn.Bind(lc.BindDN, lc.BindPassword)
		if err != nil {
			return true, user, err
		}
	}

	return true, user, nil
}

// GetGroupsOfUser returns the group for a user.
func (lc *LDAPClient) GetGroupsOfUser(username string) ([]string, error) {
	err := lc.Connect()
	if err != nil {
		return nil, err
	}

	searchRequest := ldap.NewSearchRequest(
		lc.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(lc.GroupFilter, username),
		[]string{"cn"}, // can it be something else than "cn"?
		nil,
	)
	sr, err := lc.Conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	groups := []string{}
	for _, entry := range sr.Entries {
		groups = append(groups, entry.GetAttributeValue("cn"))
	}
	return groups, nil
}
