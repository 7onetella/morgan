package ldapclient

import (
	"fmt"
	"log"
	"testing"

	"gopkg.in/ldap.v2"
)

func TestCreatePerson(t *testing.T) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "localhost", 10389))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// Add a description, and replace the mail attributes
	modify := ldap.NewModifyRequest("uid=barf,ou=people,dc=example,dc=com")
	modify.Add("description", []string{"User description"})
	// modify.Replace("mail", []string{"user@example.org"})

	err = l.Modify(modify)
	if err != nil {
		log.Fatal(err)
	}
}
