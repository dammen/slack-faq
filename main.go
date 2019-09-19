package main

import (
	//"io"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Tokens struct {
	ChannelsURL string `json: "channelsURL"`
	AppToken    string `json: "appToken"`
}

//Create a struct that holds information to be displayed in our HTML file
type Content struct {
	Name        string
	Time        string
	Topic       string
	Messages    string
	ChannelName string
}

//Go application entrypoint
func main() {
	tokens := Tokens{}
	file, _ := ioutil.ReadFile("tokens.json")
	_ = json.Unmarshal([]byte(file), &tokens)

	// The app token defines which slack workspace you are accessing
	appToken := tokens.AppToken
	count := "10"
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
	fmt.Println("Fetching messages...")
	for _, s := range channels[2:4] {
		channel := s.(map[string]interface{})
		ids = append(ids, channel["id"].(string))
		names = append(names, channel["name"].(string))
		messagesURL := "https://slack.com/api/conversations.history?token=" + appToken + "&channel=" + channel["id"].(string) + "&count=" + count
		fmt.Println("Fetching from channel " + channel["name"].(string) + "...")
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

	// Fetch messages from channel
	/*
	   messages_url := "https://slack.com/api/conversations.history?token=" + app_token + "&channel=" + ids[0]
	   resp, err = http.Get(messages_url)
	   if err != nil{
	      log.Fatalln(err)
	   }
	   defer resp.Body.Close()

	   body, err = ioutil.ReadAll(resp.Body)
	   if err != nil{
	      log.Fatalln(err)
	   }
	   var m interface{}
	   json.Unmarshal(body, &m)
	   message_data := m.(map[string]interface{})
	   messages := message_data["messages"].([]interface{})
	   var text string
	   if len(messages) > 0{
	      message := messages[0].(map[string]interface{})
	      text = message["text"].(string)
	   } else {
	      text = "No message in channel"
	   }
	*/

	jsonString, err := json.Marshal(filteredMap)

	topics := strings.Join(names, " ")
	//Instantiate a Welcome struct object and pass in some random information.
	//We shall get the name of the user as a query parameter from the URL
	content := Content{"Jonas", time.Now().Format(time.Stamp), topics, string(jsonString), names[0]}

	//We tell Go exactly where we c["an find our html file. We ask Go to parse the html file (Notice
	// the relative path). We wrap it in a call to template.Must() which handles any errors and halts if there are fatal errors

	templates := template.Must(template.ParseFiles("templates/welcome-template.html"))

	//Our HTML comes with CSS that go needs to provide when we run the app. Here we tell go to create
	// a handle that looks in the static directory, go then uses the "/static/" as a url that our
	//html can refer to when looking for our css and other files.

	http.Handle("/static/", //final url can be anything
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static")))) //Go looks in the relative "static" directory first using http.FileServer(), then matches it to a
	//url of our choice as shown in http.Handle("/static/"). This url is what we need when referencing our css files
	//once the server begins. Our html code would therefore be <link rel="stylesheet"  href="/static/stylesheet/...">
	//It is important to note the url in http.Handle can be whatever we like, so long as we are consistent.

	//This method takes in the URL path "/" and a function that takes in a response writer, and a http request.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//If errors show an internal server error message
		//I also pass the welcome struct to the welcome-template.html file.
		if err := templates.ExecuteTemplate(w, "welcome-template.html", content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

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
	writeAble, err := json.Marshal(messageData)
	if err != nil {
		log.Fatal(err)
	}
	writeToFile(name, writeAble)

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
