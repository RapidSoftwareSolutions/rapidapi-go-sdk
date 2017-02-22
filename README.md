<p align="center">
  <img src="https://storage.googleapis.com/rapid_connect_static/static/github-header.png" width=350 />
</p>

## Overview
RapidAPI is the world's first opensource API marketplace. It allows developers to discover and connect to the world's top APIs more easily and manage multiple API connections in one place.

## GO SDK

First of all, Install the package from source code

```go
go get github.com/RapidSoftwareSolutions/rapidapi-go-sdk/RapidAPISDK
```
The executable will be produced under $GOPATH/bin in your file system

## Initialization
Import RapidAPISDK by putting the following code in the head of your file
```go
import (
	"github.com/RapidSoftwareSolutions/rapidapi-go-sdk/RapidAPISDK"
)
```
  
Now initialize it using:
```go
rapidApi := RapidAPISDK.RapidAPI{"PROJECT", "TOKEN"}
```
  
## Usage

First of all, we will prepare the body, we will use a map which You can add as many arguments as you wish due to the API requirements. 
```go
params := map[string]RapidAPISDK.Param{
		"parameterName1": {"data", "value1"},
		"parameterName2": {"data", "value2"},
		"parameterName3": {"data", "value3"},
}
 ```
 Note: You have three kinds of parameters you can send - data, file or writer. A file parameter is a file path
 and a writer is a Write object (of a file).
 A data parameter is any other parameter (including url).
 
To use any block in the marketplace, just copy it's code snippet and paste it in your code. 
For example, the following is the snippet for the **MicrosoftComputerVision.analyzeImage** block:
 ```go
response := rapidApi.Call("MicrosoftComputerVision", "analyzeImage", params) 
if response["success"] != nil {
  fmt.Println(response["success"])
} else {
  fmt.Println(response["error"])
}
 ```
 
**Notice** that if you make an invalid block call (for example - the package you refer to does not exist) the program will 
exit using panic and you will see the error message there.


## Using Files
Whenever a block in RapidAPI requires a file, you can either pass a URL to the file, a path to the file or use Writer.

#### URL
The following code will call the block MicrosoftComputerVision.analyzeImage with a URL of an image:
```go
params := map[string]RapidAPISDK.Param{
		"subscriptionKey": {"data", "*****"},
		"image":           {"data", "https://i.ytimg.com/vi/opKg3fyqWt4/hqdefault.jpg"},
		"details":         {"data", ""},
		"visualFeatures":  {"data", ""},
}
response := rapidApi.Call("MicrosoftComputerVision", "analyzeImage", params)

```
#### Post File
If the file is locally stored, you can just use this line instead:
```go
"image": {"file","file/path"}
```
Or you can post a file as a Writer:
```go
file, err := os.Open("file/path")
if err != nil {
  panic(err)
}
params := map[string]RapidAPISDK.Param{
  "subscriptionKey": {"data", "*****"},
  "image":           {"writer", file},
  "details":         {"data", ""},
  "visualFeatures":  {"data", ""},
}
defer file.Close()
```

## Webhook events
You can listen to webhook events like so:

```go
	rapidApi := RapidAPISDK.RapidAPI{"PROJECT", "KEY"}
	params := map[string]string{
		"token": "slash_command_token",
		"command": "/command"}
	callbacks := make(map[string]func(msg interface{}))
	callbacks["onJoin"] = func (msg interface{}) { fmt.Println("Joined!") }
	callbacks["onMessage"] = func (msg interface{}) {
		fmt.Println("Got message!")
		fmt.Println(msg)
	}
	callbacks["onClose"] = func (msg interface{}) { fmt.Println("Closed!") }
	rapidApi.Listen("Slack", "slashCommand", params, callbacks)
```

##Issues:

As this is a pre-release version of the SDK, you may expirience bugs. Please report them in the issues section to let us know. You may use the intercom chat on rapidapi.com for support at any time.

##License:

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
