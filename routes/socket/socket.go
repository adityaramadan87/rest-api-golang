package socket

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"log"
)

func SetRoutes() *socketio.Server {

	socket, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	//socket routes
	socket.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	return socket

}
