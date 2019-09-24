package main

import (
	//"io"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var production = true
var enableWebsocket = false

// Tokens holds information necesarry to access the correct slakc-workspace
type Tokens struct {
	ChannelsURL string `json: "channelsURL"`
	AppToken    string `json: "appToken"`
}

// Content holds information to be displayed in our HTML file
type Content struct {
	Name        string
	Time        string
	Topic       string
	Messages    string
	ChannelName string
}

var available_channels = make(map[string]int)

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// We'll need to check the origin of our connection
	// this will allow us to make requests from our React
	// development server to here.
	// For now, we'll do no checking and just allow any connection
	CheckOrigin: func(r *http.Request) bool { return true },
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println("Fetching messages from the " + string(p) + " channel")
		if val, ok := available_channels[string(p)]; ok {
			if err := conn.WriteMessage(messageType, []byte("found "+strconv.Itoa(val)+" messages")); err != nil {
				log.Println(err)
				return
			}
		} else {
			if err := conn.WriteMessage(messageType, []byte("Couldn´t find a message with the name "+string(p))); err != nil {
				log.Println(err)
				return
			}

		}

	}
}

// define our WebSocket endpoint
func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//TODO: add host handle e.g. view/build
		fmt.Fprintf(w, "Simple Server")
	})
	// map our `/ws` endpoint to the `serveWs` function
	http.HandleFunc("/ws", serveWs)
}

//Go application entrypoint
func main() {
	tokens := Tokens{}
	file, _ := ioutil.ReadFile("tokens.json")
	_ = json.Unmarshal([]byte(file), &tokens)

	// The app token defines which slack workspace you are accessing
	appToken := tokens.AppToken
	count := "2"
	// fetch channels from the slack workspace
	channelsURL := tokens.ChannelsURL + appToken
	resp, err := http.Get(channelsURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var v interface{}
	json.Unmarshal(body, &v)
	data := v.(map[string]interface{})
	channels := data["channels"].([]interface{})
	var ids []string
	var names []string
	info := make(chan string)
	messages := make(chan map[string][]string, 500)
	//available_channels := new(map[string]int)
	fmt.Println("Fetching messages...")
	for _, s := range channels[:10] {
		channel := s.(map[string]interface{})
		ids = append(ids, channel["id"].(string))
		names = append(names, channel["name"].(string))
		available_channels[channel["name"].(string)] = rand.Intn(100)
		messagesURL := "https://slack.com/api/conversations.history?token=" + appToken + "&channel=" + channel["id"].(string) + "&count=" + count
		//fmt.Println("Fetching from channel " + channel["name"].(string) + "...")
		go fetch(messagesURL, channel["name"].(string), info, messages)
		fmt.Println(<-info)
	}
	var filteredMessages []string
	var dispMessages string

	filteredMap := make(map[string]string)
	close(messages)
	for elem := range messages {
		for key, value := range elem {
			for m := range value {
				if len(value[m]) > 0 {
					if value[m][0] != []byte("<")[0] {
						if strings.Contains(value[m], "?") {
							//fmt.Println(elem[m])
							filteredMessages = append(filteredMessages, value[m])
						}
					}
				}
			}
			if len(filteredMessages) > 0 {
				dispMessages = strings.Join(filteredMessages, "/\n/")
				filteredMap[key] = dispMessages
				//writeToFile(key, filteredMessages)
			}
			filteredMessages = nil
		}
	}

	fmt.Printf("Fetched %d channel ids \n", len(ids))
	if production {
		buildHandler := http.FileServer(http.Dir("./view/build"))
		http.Handle("/", buildHandler)
		http.Handle("/*/", buildHandler)
	}
	if enableWebsocket {
		setupRoutes()
	}

	//Start the web server, set the port to listen to 8080. Without a path it assumes localhost
	//Print any errors from starting the webserver using fmt
	fmt.Println("Listening")
	fmt.Println(http.ListenAndServe(":8080", nil))
}

func fetch(url string, name string, info chan<- string, messages chan<- map[string][]string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		info <- fmt.Sprint(err)
		return
	}
	//nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()
	if err != nil {
		info <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var m interface{}
	json.Unmarshal(body, &m)
	messageData := m.(map[string]interface{})
	messageData["topic"] = name
	//writeAble, err := json.Marshal(messageData)
	if err != nil {
		log.Fatal(err)
	}
	//writeToFile(name, writeAble)

	//tempMessages := messageData["messages"].([]interface{})
	nmessages := len(messageData["messages"].([]interface{}))
	var messagesAsText []string
	messagesMap := make(map[string][]string)

	if nmessages == 0 {
		fmt.Println("No messages in channel " + name)
		messagesAsText = append(messagesAsText, "")

	} else {
		for _, text := range messageData["messages"].([]interface{}) {
			t := text.(map[string]interface{})
			fmt.Println("------------------------------------------------------------------------------------------------------")
			fmt.Println(t)
			te := t["text"].(string)
			messagesAsText = append(messagesAsText, te)
		}
	}
	messagesMap[name] = messagesAsText
	messages <- messagesMap

	secs := time.Since(start).Seconds()
	info <- fmt.Sprintf("%.2fs %7d %s \n", secs, nmessages, name)
}

func writeToFile(topic string, messages []byte) {
	// If the file doesn't exist, create it, or append to the file
	//messagesAsMap := make(map[string][]string)
	//messagesAsMap[topic] = messages
	//messagesAsJSON, err := json.Marshal(messagesAsMap)
	//if err != nil {
	//	log.Fatal(err)
	//}
	f, err := os.OpenFile(topic+".json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(messages)); err != nil {
		f.Close() // ignore error; Write error takes precedence
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
