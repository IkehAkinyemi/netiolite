package netiolite

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	conn, err := ln.Accept()
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	buf := make([]byte, 164)
	n, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}

	fmt.Println(n)

	for {
		n, err = conn.Read(buf[:])
		fmt.Println(n)
		if err != nil {
			break
		}
	}
	fmt.Println(n)

	// file, err := ln.(*net.TCPListener).File()
	// if err != nil {
	// 	panic(err)
	// }

	// epoll := newPoll()

	// for {
	// 	fds := epoll.wait(-1)
	// 	fmt.Println(fds)
	// }
}