package utils

import "github.com/nadezhda/documents/services/base"

const (
	BACKEND_ENV = "EMPLOYEE_BACKEND_ADDRESS"
)

func GetBackendAddress() string {
	address := base.GetEnv(BACKEND_ENV, "http://localhost:7074")
	return address
}
