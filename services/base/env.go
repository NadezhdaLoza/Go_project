package base

	import (

		"os"
	)

    func GetEnv(key string, defaultValue string)string{
        if value := os.Getenv(key); value ==""{
            return defaultValue
        }else{
            return value
        }
    }


