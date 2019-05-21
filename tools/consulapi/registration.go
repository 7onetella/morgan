package consulapi

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
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/hashicorp/consul/api"
)

// Register registers this docker instance with local consul
// excert from https://github.com/fabiolb/fabio/blob/master/demo/server/server.go
func Register(name, host string, port int, checkpath string, tags []string, clientaddr string) {

	addr := fmt.Sprintf("%s:%d", host, port)

	checkScheme := "http"
	check := &api.AgentServiceCheck{
		HTTP:     checkScheme + "://" + addr + checkpath,
		Interval: "10s",
		Timeout:  "1s",
	}

	// register service with health check
	serviceID := name + "-" + addr
	service := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    name,
		Port:    port,
		Address: host,
		Tags:    tags,
		Check:   check,
	}

	config := &api.Config{Address: clientaddr, Scheme: "http"}
	client, err := api.NewClient(config)
	if err != nil {
		log.Println(err)
		return
	}

	if err := client.Agent().ServiceRegister(service); err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("Registering %s service %q in consul with tags %q", "http", name, strings.Join(tags, ","))

	// run until we get a signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-quit

	// deregister service
	if err := Deregister(serviceID, clientaddr); err != nil {
		log.Printf("Deregistering service %q in consul failed - %v", name, err)
		return
	}
	log.Printf("Deregistering service %q in consul", name)
}

// Deregister deregisters consul service instance
func Deregister(serviceID, clientaddr string) error {

	config := &api.Config{Address: clientaddr, Scheme: "http"}
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}
	// deregister service
	if err := client.Agent().ServiceDeregister(serviceID); err != nil {
		return err
	}

	return nil

}
