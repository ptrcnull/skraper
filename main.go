package main

import (
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"log"
	"sync"
)

func client(wg *sync.WaitGroup) {
	c, err := gosocketio.Dial(
		gosocketio.GetUrl("skribbl.io", 5001, true),
		transport.GetDefaultWebsocketTransport(),
	)
	if err != nil {
		panic(err)
	}
	_ = c.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("O KURWA")
		_ = c.Emit("userData", UserData{
			Name:          "rabit",
			Code:          "",
			Avatar:        []int{17, 21, 16, -1},
			Join:          "",
			Language:      "English",
			CreatePrivate: false,
		})
	})

	c.On("lobbyConnected", func(c *gosocketio.Channel, iargs interface{}) {
		log.Println("X kurwa D")
		args, ok := iargs.(map[string]interface{})
		if !ok {
			log.Println("kurwa zjebalo sie cos")
			return
		}
		log.Println(args)
		log.Println(args["name"])
		log.Println(args["myID"])
	})

	c.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Printf("%s got fucked.\n", c.Id())
		wg.Done()
	})
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go client(&wg)
	wg.Wait()
}