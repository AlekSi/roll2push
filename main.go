package main

import (
	"encoding/json"
	"flag"
	"github.com/AlekSi/pushover"
	"github.com/AlekSi/roll2push/rollbar"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	User string
)

func hook(rw http.ResponseWriter, req *http.Request) {
	l := log.New(os.Stderr, req.RemoteAddr+" at ", log.Ldate|log.Lmicroseconds)
	var err error
	defer func() {
		if err == nil {
			rw.WriteHeader(204)
		} else {
			l.Print(err)
			rw.WriteHeader(500)
		}
	}()

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}

	var event rollbar.Event
	err = json.Unmarshal(b, &event)
	if err != nil {
		return
	}

	switch event.EventName {
	case "test":
		var test string
		err = json.Unmarshal(event.Data["message"], &test)
		if err != nil {
			return
		}

		l.Print(test)
		l.Printf("pushover: %v", pushover.SendMessage(User, "test", test))

	case "deploy":
		l.Printf("event is not handled yet: %s", event.EventName)

	default:
		if strings.HasSuffix(event.EventName, "_item") {
			var item rollbar.ItemEvent
			err = json.Unmarshal(event.Data["item"], &item)
			if err != nil {
				return
			}

			l.Printf("%+v", item)
			l.Printf("pushover: %v", pushover.SendMessage(User, item.Environment, item.Title))

		} else {
			l.Printf("unexpected event: %s", event.EventName)
		}
	}

	return
}

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)

	token := flag.String("pushover-token", "", "")
	flag.StringVar(&User, "pushover-user", "", "")
	flag.Parse()

	pushover.DefaultClient.Token = *token

	http.HandleFunc("/", hook)

	addr := ":8080"
	log.Printf("Waiting for Rollbar events on %s.", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
