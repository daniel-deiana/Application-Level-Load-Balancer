package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func main() {
	client := &http.Client {}	

	resp, err := client.Get("http://localhost:8080/lb")
	
	if err != nil {
		fmt.Println("Error code is: ", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Response is: ", resp)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("The response body is: ", string(body))	
}
