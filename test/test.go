package main

import (
	"fmt"
	"os"
	"../RapidAPISDK"
)

func handleResponse(response map[string]interface{}){
	if response["success"] != nil {
		fmt.Println(response["success"])
	} else {
		fmt.Println(response["error"])
	}
}

func TestPublicPack(rapidApi RapidAPISDK.RapidAPI) {
	params := map[string]RapidAPISDK.Param{
		"apiKey": {"data","AIzaSyCDogEcpeA84USVXMS471PDt3zsG-caYDM"},
		"string": {"data", "test"},
		"targetLanguage": {"data", "he"},
		"sourceLanguage": {"data",""},
	}
	response := rapidApi.Call("GoogleTranslate", "translate", params)
	handleResponse(response)
}

func TestPackWithImage(rapidApi RapidAPISDK.RapidAPI) {
	params :=  map[string]RapidAPISDK.Param{
		"subscriptionKey": {"data", "57e9164516844d99ae455a9953aca0c2"},
		"image" : {"file","test/cute_dog.jpg" },
		"details": {"data", ""},
		"visualFeatures": {"data",""},
	}
	response := rapidApi.Call("MicrosoftComputerVision", "analyzeImage", params)
	handleResponse(response)
}

func TestPackWithWriter(rapidApi RapidAPISDK.RapidAPI) {

	file, err := os.Open("test/cute_dog.jpg")
	if err != nil {
		panic(err)
	}
	params := map[string]RapidAPISDK.Param{
		"subscriptionKey": {"data", "57e9164516844d99ae455a9953aca0c2"},
		"image":           {"writer", file},
		"details":         {"data", ""},
		"visualFeatures":  {"data", ""},
	}
	defer file.Close()

	response := rapidApi.Call("MicrosoftComputerVision", "analyzeImage", params)
	handleResponse(response)
}

func TestPackWithURL(rapidApi RapidAPISDK.RapidAPI) {

	params := map[string]RapidAPISDK.Param{
		"subscriptionKey": {"data", "57e9164516844d99ae455a9953aca0c2"},
		"image":           {"data", "https://i.ytimg.com/vi/opKg3fyqWt4/hqdefault.jpg"},
		"details":         {"data", ""},
		"visualFeatures":  {"data", ""},
	}

	response := rapidApi.Call("MicrosoftComputerVision", "analyzeImage", params)
	handleResponse(response)
}

func TestListen(rapidApi RapidAPISDK.RapidAPI) {
	params := map[string]string{
		"command": "/send_test",
		"token": "ydt3vFyVEoW51ZFCC2i5QKab",
	}

	params2 := map[string]string{
		"command": "/different_test",
		"token": "et7cSI58sSuHnSeSa4DQP8hn",
	}

	on_join := make(chan bool)
	on_message := make(chan interface{})
	on_error := make(chan interface{})
	on_close := make(chan interface{})
	
	go rapidApi.Listen("Slack", "slashCommand", params, on_join, on_message, on_error, on_close)
	go rapidApi.Listen("Slack", "slashCommand", params2, on_join, on_message, on_error, on_close)

	for {
		select {
		case <-on_join:
			fmt.Println("Joined")
		case message := <-on_message:
			fmt.Println(message)
		case <-on_close:
			fmt.Println("closed")
		case err := <-on_error:
			fmt.Println("error")
			fmt.Println(err)
		}
	}
}

func main() {
	// rapidApi := RapidAPISDK.RapidAPI{"withoutImage", "72352b8b-9384-4a9a-abb1-195d5e234418"}

	rapidApi := RapidAPISDK.RapidAPI{"Dashboard", "0b7f82c1-cf1d-4e02-af9a-de129afe54b2"}
	TestListen(rapidApi, done)

	TestPublicPack(rapidApi)
	TestPackWithImage(rapidApi)
	TestPackWithURL(rapidApi)
	TestPackWithWriter(rapidApi)
}
