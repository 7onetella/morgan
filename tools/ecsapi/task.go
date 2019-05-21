package ecsapi

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
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// GetContainerName returns the container name
func GetContainerName() (string, error) {
	container, err := getContainer()
	if err != nil {
		return "", err
	}

	return container.Name, nil
}

// GetHostPort retrieves ecs task host port
func GetHostPort() (int, error) {
	container, err := getContainer()
	if err != nil {
		return 0, err
	}

	return container.Ports[0].HostPort, nil
}

// GetHost hits http://172.17.0.1:51678/v1/tasks?dockerid=
func GetHost() (string, error) {
	data, err := getURL("http://169.254.169.254/latest/meta-data/local-ipv4")
	if err != nil {
		return "", err
	}
	return string(data), err
}

func dockerid() (string, error) {
	data, err := readFile("/proc/1/cpuset")
	if err != nil {
		return "", err
	}

	s := string(data)

	// sample /ecs/973231eb-92ce-48a9-8eca-f28abcf0f7bc/37fb340c39aff98d4d111f92d1653b50342f19456fd6dcd65dac8b11a5234916
	// last chunk is what we want
	dockerid := strings.Split(s, "/")[3]

	// dockerid := strings.Replace(string(data), "/docker/", "", 1)

	// last character has a carriage return
	return strings.TrimSpace(dockerid), err
}

func readFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

func getURL(url string) ([]byte, error) {
	client := http.Client{
		Timeout: time.Second * 3,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	return data, err
}

// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-agent-introspection.html
func taskMetaDataJSON(dockerid string) (taskJSON, error) {
	t := taskJSON{}
	// loop until "KnownStatus":"PENDING", is RUNNING
	data, err := getURL("http://172.17.0.1:51678/v1/tasks?dockerid=" + dockerid)
	if err != nil {
		return t, err
	}
	log.Println("task: ", string(data))
	json.Unmarshal(data, &t)
	return t, nil
}

func getTask() (taskJSON, error) {

	task := taskJSON{}

	dockerID, err := dockerid()
	log.Printf("dockerid: [%s]\n", dockerID)
	if err != nil {
		return task, err
	}

	for i := 0; i < 30; i++ {
		// query ecs agent for the given docker instance
		t, err := taskMetaDataJSON(dockerID)
		if err != nil {
			return task, err
		}

		if t.KnownStatus != "RUNNING" {
			time.Sleep(time.Second * 10)
			continue
		}

		task = t
		break
	}

	return task, nil
}

// Container refers to docker container
type Container struct {
	DockerID   string `json:"DockerId"`
	DockerName string `json:"DockerName"`
	Name       string `json:"Name"`
	Ports      []struct {
		ContainerPort int    `json:"ContainerPort"`
		Protocol      string `json:"Protocol"`
		HostPort      int    `json:"HostPort"`
	} `json:"Ports"`
}

// TaskJSON ecs task json
type taskJSON struct {
	Arn           string      `json:"Arn"`
	DesiredStatus string      `json:"DesiredStatus"`
	KnownStatus   string      `json:"KnownStatus"`
	Family        string      `json:"Family"`
	Version       string      `json:"Version"`
	Containers    []Container `json:"Containers"`
}

func getContainer() (Container, error) {
	emptyContainer := Container{}

	task, err := getTask()
	if err != nil {
		return emptyContainer, err
	}

	// this shouldn't happen if it did, let's avoid runtime error
	if len(task.Containers) == 0 {
		return emptyContainer, errors.New("task data empty")
	}

	return task.Containers[0], nil
}
