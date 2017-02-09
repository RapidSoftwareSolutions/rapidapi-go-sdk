package RapidAPISDK

import (
	"io"
	"os"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"mime/multipart"
)

// base URL for block calls
const baseUrl string = "https://rapidapi.io/connect"

/*
 * Create rapidAPI connect type
 */
type RapidAPI struct {
	Project, Key string
}

/*
 * Represents a parameter value and its type
 */
type Param struct {
	Type  string
	Value interface{}
}

/**
 * Call a block
 * @param pack Package of the block
 * @param block Name of the block
 * @param body Arguments to send to the block
 * @returns {map[string]interface{}} response
 */
func (rapidApi RapidAPI) Call(pack, block string, params map[string]Param) map[string]interface{} {

	body, writer := createBody(params)
	url := blockURLBuilder(pack, block)
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.SetBasicAuth(rapidApi.Project, rapidApi.Key)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bodyRes, _ := ioutil.ReadAll(resp.Body)

	return renderResponse(bodyRes)

}

/*
 * Build a URL for a block call
 * @param pack Package where the block is
 * @param block Block to be called
 * @returns {string} Generated URL
 */
func blockURLBuilder(pack string, block string) string {
	return baseUrl + "/" + pack + "/" + block
}

/*
 * Creates the request body
 * @param params body params to parse
 * @returns {io.Reader, multipart.Writer} for the http request
 */
func createBody(params map[string]Param) (io.Reader, multipart.Writer) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// add all params
	for key, val := range params {
		switch val.Type {
		case "data":
			_ = writer.WriteField(key, val.Value.(string))
		case "file":
			// add file parameter
			file, err := os.Open(val.Value.(string))
			if err != nil {
				panic(err)
			}
			//defer file.Close()
			part, err := writer.CreateFormFile(key, "file")
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(part, file)
			file.Close()

		case "writer":
			part, err := writer.CreateFormFile(key, "file")
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(part, val.Value.(*os.File))
			val.Value.(*os.File).Close()
		default:
			panic(val.Type + " is not a valid parameter type")
		}
	}

	err := writer.Close()
	if err != nil {
		panic(err)
	}
	return body, *writer

}

/*
 * Render the response for the user
  * @param bodyRes Body response
  * @returns {map[string]interface{}} rendered response
*/
func renderResponse(bodyRes []byte) map[string]interface{} {
	var value interface{}
	err := json.Unmarshal(bodyRes, &value)
	if err != nil {
		panic(string(bodyRes))
	}
	var outcome = make(map[string]interface{})
	mapRes := value.(map[string]interface{})

	if (mapRes["outcome"]).(string) == "success" {
		outcome["success"] = mapRes["payload"]
	} else {
		outcome["error"] = mapRes["payload"]

	}
	return outcome

}
