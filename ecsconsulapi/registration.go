package ecsconsulapi

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

	"github.com/7onetella/morgan/tools/consulapi"
	"github.com/7onetella/morgan/tools/ecsapi"
)

// Registrator config for registering with Consul
type Registrator struct {
	Name        string // if not specified NAME environment variable is used
	Healthcheck string // if not specified /health is used
	URLPrefix   string // if not specified URLPREFIX environment variable is used
	ClientAddr  string // if not specified 169.254.1.1:8500 is used
}

// New initialize with health check only
func New(healthcheck string) Registrator {
	return Registrator{Healthcheck: healthcheck}
}

// ServiceRegister registers service with consule
func (r Registrator) ServiceRegister() {

	if len(r.Name) == 0 {
		r.Name = os.Getenv("NAME")
	}

	host, _ := ecsapi.GetHost()

	port, _ := ecsapi.GetHostPort()

	if len(r.Healthcheck) == 0 {
		r.Healthcheck = "/health"
	}

	if len(r.URLPrefix) == 0 {
		r.URLPrefix = os.Getenv("URLPREFIX")
	}

	if len(r.ClientAddr) == 0 {
		r.ClientAddr = "169.254.1.1:8500"
	}

	tags := []string{"urlprefix-" + os.Getenv("URLPREFIX")}

	consulapi.Register(r.Name, host, port, r.Healthcheck, tags, r.ClientAddr)
}
