package main

import (
	"fmt"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"log"
	"sync"
)

func client(wg *sync.WaitGroup, e *Chans, port int, workerID int, key string) {
	wg.Add(1)

	name := "daddy"
	if key != "" {
		name = "child"
	}

	prt := func(args ...interface{}) {
		log.Printf("%d[%d][%s]: %s", port, workerID, name, fmt.Sprint(args...))
	}

	c, err := gosocketio.Dial(
		gosocketio.GetUrl("skribbl.io", port, true),
		transport.GetDefaultWebsocketTransport(),
	)
	if err != nil {
		panic(err)
	}

	emit := func(method string, data interface{}) {
		if err := c.Emit(method, data); err != nil {
			prt(fmt.Sprintf("Error while sending %s:", method))
			prt(err)
			wg.Done()
			c.Close()
		}
	}

	_ = c.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		prt("O KURWA")

		createPrivate := false
		if key == "" {
			createPrivate = true
		}

		emit("userData", UserData{
			Name:          "rabit",
			Code:          "",
			Avatar:        []int{17, 21, 16, -1},
			Join:          key,
			Language:      "English",
			CreatePrivate: createPrivate,
		})
	})

	var myID float64

	_ = c.On("lobbyConnected", func(c *gosocketio.Channel, iargs interface{}) {
		args, ok := iargs.(map[string]interface{})
		if !ok {
			prt("Error when casting lobby data to map[string]")
			prt(iargs)
			return
		}
		prt("Connected. Key: ", args["key"], ", ID: ", args["myID"])
		myID = args["myID"].(float64)
		e.Key<-args["key"].(string)
	})


	_ = c.On("lobbyPlayerConnected", func(c *gosocketio.Channel) {
		prt("lobbyPlayerConnected")
		if key == "" {
			prt("Starting!")
			emit("lobbyGameStart", "")
		}
	})

	_ = c.On("lobbyLobby", func(c *gosocketio.Channel) {
		prt("lobbyLobby")
		if key == "" {
			prt("Starting!")
			emit("lobbyGameStart", "")
		}
	})

	_ = c.On("lobbyChooseWord", func(c *gosocketio.Channel, ev WordsEvent) {
		prt("O KURWA MAM SŁOWA")
		prt(ev)
		if len(ev.Words) != 0 {
			prt(fmt.Sprintf("%v == %v: %v (%v)", ev.Id, myID, ev.Words, len(ev.Words)))
			emit("lobbyChooseWord", 0)
		}
	})

	_ = c.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		prt("got fucked.")
		wg.Done()
	})

	go func() {
		for word := range e.Word {
			if !c.IsAlive() {
				continue
			}
			emit("chat", word)
		}
	}()
}

func create_room(port int, id int, mainwg *sync.WaitGroup) {
	mainwg.Add(1)
	defer mainwg.Done()

	var wg sync.WaitGroup

	events := Chans{
		Word: make(chan string),
		Key:  make(chan string),
		NewWord: make(chan string),
	}

	go client(&wg, &events, port, id, "")

	key := <-events.Key

	go client(&wg, &events, port, id, key)

	wg.Wait()
}

func main() {
	var wg sync.WaitGroup

	for port := 1; port <= 4; port++ {
		for id := 0; id < 3; id++ {
			go create_room(5000 + port, id, &wg)
		}
	}

	wg.Wait()
}
