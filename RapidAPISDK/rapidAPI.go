package RapidAPISDK

import (
	"io"
	"os"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"mime/multipart"
	"github.com/gorilla/websocket"
	"time"
)

// base URL for block calls
const baseUrl string = "https://rapidapi.io/connect"
const callbackBaseUrl string = "https://webhooks.rapidapi.com"
const websocketBaseUrl string = "wss://webhooks.rapidapi.com"

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

func getTokenUrl(user_id string) string {
	return callbackBaseUrl + "/api/get_token?user_id=" + user_id
}

func socketUrl(token string) string {
	return websocketBaseUrl + "/socket/websocket?token=" + token
}

func getToken(user_id string, rapidApi RapidAPI) string {
	url := getTokenUrl(user_id)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(rapidApi.Project, rapidApi.Key)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bodyRes, _ := ioutil.ReadAll(resp.Body)
	var value interface{}
	if (json.Unmarshal(bodyRes, &value) != nil) {
		panic(string(bodyRes))
	}
	mapRes := value.(map[string]interface{})
	return mapRes["token"].(string)
}

func (rapidApi RapidAPI) Listen(pack string, event string, params map[string]string,
	on_join chan bool, on_message chan interface{},
	on_error chan interface{}, on_close chan interface{}) {
	user_id := pack + "." + event + "_" + rapidApi.Project + ":" + rapidApi.Key
	token := getToken(user_id, rapidApi)
	sock_url := socketUrl(token)
	c, _, err := websocket.DefaultDialer.Dial(sock_url, nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	toclose := make(chan struct{})
	defer close(toclose)

	connect := make(map[string]interface{})
	connect["topic"] = "users_socket:"+token
	connect["event"] = "phx_join"
	connect["ref"] = "1"
	connect["payload"] = params

	jsoned, _ := json.Marshal(connect)
	if err := c.WriteMessage(websocket.TextMessage, []byte(jsoned)); err != nil {
	 	panic(err)
	}

	var payload map[string]interface{}
	go func(conn *websocket.Conn) {
		heartbeat := map[string]interface{} {
			"topic": "phoenix",
			"event": "heartbeat",
			"ref": "1",
			"payload": make(map[string]string),
		}
		jsoned_heartbeat, _ := json.Marshal(heartbeat)
		for range time.Tick(30*time.Second) {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(jsoned_heartbeat)); err != nil {
				panic(err)
			}
		}
	}(c)

	for {
		_, message, err := c.ReadMessage();
		if err != nil {
			on_error <- true
			return
		}
		if err := json.Unmarshal(message, &payload); err != nil {
			on_error <- err
			panic(err)
		}
		msg := payload["payload"].(map[string]interface{})
		if payload["event"].(string) == "joined" {
			on_join <- true
		} else if payload["event"].(string) == "new_msg" && msg["token"] == nil {
			on_error <- msg["body"]
		} else if payload["event"].(string) == "new_msg" {
			on_message <- msg["body"]
		}
	}
	on_close <- true
	return

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
