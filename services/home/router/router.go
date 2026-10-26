package router

	import (
		"net/http"
		"net/http/httputil"
		"net/url"

		"github.com/gorilla/mux"
		"github.com/nadezhda/documents/services/base"
	)


	func handleWithProxy(router *mux.Router, path string, address string) {
		router.PathPrefix(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			url, _ := url.Parse(address)
			proxy := httputil.NewSingleHostReverseProxy(url)
			proxy.ServeHTTP(w, r)
		})
	}

	func GetRouter() http.Handler {
		router := mux.NewRouter()

		// Получение адреса товаров из переменной окружения или же значение по умолчанию
		productAddress := base.GetEnv("PRODUCT_ADDRESS", "http://localhost:8081")
		handleWithProxy(router, "/product", productAddress)

		organizationAddress := base.GetEnv("ORGANIZATION_ADDRESS", "http://localhost:8082")
		handleWithProxy(router, "/organization", organizationAddress)

		accountAddress := base.GetEnv("ACCOUNT_ADDRESS", "http://localhost:8083")
		handleWithProxy(router, "/account", accountAddress)

		employeeAddress := base.GetEnv("EMPLOYEE_ADDRESS", "http://localhost:8084")
		handleWithProxy(router, "/employee", employeeAddress)

		customerAddress := base.GetEnv("CUSTOMER_ADDRESS", "http://localhost:8085")
		handleWithProxy(router, "/customer", customerAddress)

		departmentAddress := base.GetEnv("DEPARTMENT_ADDRESS", "http://localhost:8088")
		handleWithProxy(router, "/department", departmentAddress)

		workAddress := base.GetEnv("WORK_ADDRESS", "http://localhost:8089")
		handleWithProxy(router, "/work", workAddress)

		// Прокси для документов
		proxyAddress := base.GetEnv("PROXY_ADDRESS", "http://localhost:8086")
		handleWithProxy(router, "/documents/proxies", proxyAddress)

		orderAddress := base.GetEnv("ORDER_ADDRESS", "http://localhost:8087")
		handleWithProxy(router, "/documents/orders", orderAddress)
		
		// ВАЖНО: Оставляем проксирование статических файлов для микросервисов Documents
		documentsStatic := base.GetEnv("DOCUMENTS_STATIC", "http://localhost:8086")
		handleWithProxy(router, "/documents/static", documentsStatic)
		
		// Общие статические файлы через прокси (чтобы микросервисы Documents могли их получить)
		commonStatic := base.GetEnv("COMMON_STATIC", "http://localhost:8086")
		handleWithProxy(router, "/static", commonStatic)
		
		// ДУБЛИРУЮЩАЯ ЗАПИСЬ для локального обслуживания статики
		router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", 
			http.FileServer(http.Dir("../../../static"))))

		return router
	}