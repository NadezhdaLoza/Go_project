package utils

import "github.com/nadezhda/documents/services/base"

const (
    EMPLOYEES_ENV     = "EMPLOYEE_BACKEND_ADDRESS"
    CUSTOMERS_ENV     = "CUSTOMER_BACKEND_ADDRESS"
    ORDERS_ENV        = "ORDER_BACKEND_ADDRESS"
    DEPARTMENTS_ENV   = "DEPARTMENT_BACKEND_ADDRESS"
    WORKS_ENV         = "WORK_BACKEND_ADDRESS"
    ORGANIZATIONS_ENV = "ORGANIZATION_BACKEND_ADDRESS"  
    PRODUCTS_ENV      = "PRODUCT_BACKEND_ADDRESS"       
)

func GetOrdersAddress() string {
    address := base.GetEnv(ORDERS_ENV, "http://localhost:7077")
    return address
}

func GetEmployeesAddress() string {
    address := base.GetEnv(EMPLOYEES_ENV, "http://localhost:7074")
    return address
}

func GetCustomersAddress() string {
    address := base.GetEnv(CUSTOMERS_ENV, "http://localhost:7075")
    return address
}

func GetDepartmentsAddress() string {
    address := base.GetEnv(DEPARTMENTS_ENV, "http://localhost:7078")
    return address
}

func GetWorksAddress() string {
    address := base.GetEnv(WORKS_ENV, "http://localhost:7079")
    return address
}

// Добавленные функции
func GetOrganizationsAddress() string {
    address := base.GetEnv(ORGANIZATIONS_ENV, "http://localhost:7072")
    return address
}

func GetProductsAddress() string {
    address := base.GetEnv(PRODUCTS_ENV, "http://localhost:7071")
    return address
}