package cmd

import (
	"fmt"
	"testing"

	"github.com/7onetella/morgan/internal/execw"
)

func TestECSCreate(t *testing.T) {

	stdout, stderr, err := execw.Exec([]string{"morgan", "aws", "ecs", "create-service", "bar-svc", "xsmall", "8080", "nginx:latest", "--cluster", "Shared",
		"--log", "debug"})

	if err != nil {
		t.Log(stdout)
		t.Log(stderr)
		t.Fatalf("morgan aws ecs create-service failed: %v", err)
	}

	fmt.Println(stdout)
	fmt.Println(stderr)

}
