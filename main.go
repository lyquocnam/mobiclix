package main

import "fmt"

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{
		"status":  status,
		"message": message,
	}
}

func main() {
	m := Message(true, "abc")
	m["account"] = "account"
	fmt.Println(m)
	//app.Run()
}
