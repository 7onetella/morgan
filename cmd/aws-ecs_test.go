package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/7onetella/morgan/internal/execw"
	"github.com/google/uuid"
)

func ServiceUUID() string {
	id := uuid.New().String()
	i := strings.LastIndex(id, "-")
	// service name length has a limit of 32 characters
	idPartial := id[i+1:]

	return "svc-" + idPartial
}

func DeleteService(service string) {
	cmd := fmt.Sprintf("morgan aws ecs delete %s", service)
	execw.Exec(strings.Split(cmd, " "))
}

func CreateService() (string, string, error) {
	service := ServiceUUID()
	cmd := fmt.Sprintf("morgan aws ecs create-service %s small 8080 nginx:latest --cluster Shared --log debug", service)

	stdout, _, err := execw.Exec(strings.Split(cmd, " "))

	return service, stdout, err
}

func LogFail(stdout, msgformat string, err error, t *testing.T) {
	if err != nil {
		t.Log(stdout)
		t.Logf(msgformat, err)
		t.Fail()
	}
}

func TestCreateService(t *testing.T) {

	service, stdout, err := CreateService()
	LogFail(stdout, "morgan aws ecs create-service failed: %v", err, t)

	fmt.Println(stdout)

	DeleteService(service)
}

func TestStopService(t *testing.T) {

	service, stdout, err := CreateService()
	LogFail(stdout, "morgan aws ecs create-service failed: %v", err, t)
	fmt.Println(stdout)

	cmd := fmt.Sprintf("morgan aws ecs stop-service %s --cluster Shared --log debug", service)
	stdout, _, err = execw.Exec(strings.Split(cmd, " "))
	LogFail(stdout, "morgan aws ecs stop-service failed: %v", err, t)

	DeleteService(service)

}

func TestStartService(t *testing.T) {

	service, stdout, err := CreateService()
	LogFail(stdout, "morgan aws ecs create-service failed: %v", err, t)
	fmt.Println(stdout)

	cmd := fmt.Sprintf("morgan aws ecs stop-service %s --cluster Shared --log debug", service)
	stdout, _, err = execw.Exec(strings.Split(cmd, " "))
	LogFail(stdout, "morgan aws ecs stop-service failed: %v", err, t)

	cmd = fmt.Sprintf("morgan aws ecs start-service %s --cluster Shared --log debug", service)
	stdout, _, err = execw.Exec(strings.Split(cmd, " "))
	LogFail(stdout, "morgan aws ecs start-service failed: %v", err, t)

	DeleteService(service)

}

func TestUpdateService(t *testing.T) {

	service, stdout, err := CreateService()
	LogFail(stdout, "morgan aws ecs create-service failed: %v", err, t)
	fmt.Println(stdout)

	cmd := fmt.Sprintf("morgan aws ecs update-service %s latest --cluster Shared --log debug", service)
	stdout, _, err = execw.Exec(strings.Split(cmd, " "))
	LogFail(stdout, "morgan aws ecs update-service failed: %v", err, t)

	DeleteService(service)

}
