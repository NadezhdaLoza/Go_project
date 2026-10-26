package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nadezhda/documents/services/base"
	"github.com/nadezhda/documents/services/proxy/frontend/router"
)

func main() {

	handler := router.GetRouter()

	fmt.Println("Listening on 8086...")
	if err := http.ListenAndServe(":8086", handler); err != nil {
		base.LogError(err)
		os.Exit(-1)
	}
}
