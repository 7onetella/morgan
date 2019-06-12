// Copyright © 2019 Seven OneTella<7onetella@gmail.com>
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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var certCmdVirtualHost string

// certCmd represents the cert command
var certCmd = &cobra.Command{
	Use:     "cert <addr>",
	Short:   "Displays cert info",
	Long:    `Displays cert info`,
	Example: "cert example.com:443",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		addrS := args[0]
		addrS = strings.TrimSpace(addrS)

		var host string
		port := "443" // assume default unless explicitly specified
		tokens := strings.Split(addrS, ":")

		host = addrS

		if len(tokens) == 2 {
			host = tokens[0]
			port = tokens[1]
		}

		certs, _, err := getCerts(host, port, certCmdVirtualHost)
		ExitOnError(err, "retrieving certs")

		if len(certs) > 0 {
			fmt.Println()
		}

		for _, c := range certs {
			fmt.Println("Issuer      : ", c.Issuer.CommonName)
			fmt.Println("Common Name : ", c.Subject.CommonName)
			fmt.Println("SANs        : ", c.DNSNames)
			fmt.Println("Not Before  : ", c.NotBefore.In(time.Local).String())
			fmt.Println("Not After   : ", c.NotAfter.In(time.Local).String())
			fmt.Println()
		}

	},
}

func init() {

	rootCmd.AddCommand(certCmd)

	flags := certCmd.Flags()

	flags.StringVar(&certCmdVirtualHost, "host", "", "host")

}

// code snippet from https://github.com/genkiroid/cert/blob/master/cert.go
// with minor modification
func getCerts(host, port, virtualhost string) ([]*x509.Certificate, string, error) {
	d := &net.Dialer{
		Timeout: time.Duration(3) * time.Second,
	}

	cfg := &tls.Config{
		InsecureSkipVerify: true,
	}

	if len(virtualhost) > 0 {
		cfg.ServerName = virtualhost
	}

	conn, err := tls.DialWithDialer(d, "tcp", host+":"+port, cfg)

	if err != nil {
		return []*x509.Certificate{&x509.Certificate{}}, "", err
	}
	defer conn.Close()

	addr := conn.RemoteAddr()
	ip, _, _ := net.SplitHostPort(addr.String())
	cert := conn.ConnectionState().PeerCertificates

	return cert, ip, nil
}
