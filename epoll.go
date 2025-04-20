// +build linux

package netiolite

import (
	"fmt"
	"net"
	"sync"
	"syscall"
	"time"
)

const (
	EPOLLET        = 0x80000000
	EPOLLEXCLUSIVE = 0x10000000
)

type poll struct {
	fd int
	events []syscall.EpollEvent
	evfds []int32
}

func newPoll(flags int) (*poll, error) {
	fd, err := syscall.EpollCreate1(flags)
	if err != nil {
		return nil, err
	}

	p := new(poll)
	p.fd = fd

	p.events = make([]syscall.EpollEvent, 64)
	p.evfds = make([]int32, len(p.events))

	return p, nil
}

func (p *poll) wait(msec time.Duration) ([]int32, error) {
	var err error
	var n int
	if msec >= 0 {
		n, err = syscall.EpollWait(p.fd, p.events, int(msec / time.Millisecond))
	} else {
		n, err = syscall.EpollWait(p.fd, p.events, -1)
	}

	if err != nil && err != syscall.EINTR {
		return nil, err
	}

	p.evfds = p.evfds[:0]
	for i := 0; i < n; i++ {
		p.evfds = append(p.evfds, p.events[i].Fd)
	}

	return p.evfds, nil
}

func (p *poll) addEvents(fd int32, flags int) error {
	if err := syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_ADD, int(fd), &syscall.EpollEvent{
		Fd: fd,
		Events: uint32(flags),
	}); err != nil {
		return err
	}

	return nil
}

func (p *poll) modEvents(fd int32, flags int) error {
	if err := syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_MOD, int(fd), &syscall.EpollEvent{
		Fd: fd,
		Events: uint32(flags),
	}); err != nil {
	return err
	}

	return nil
}

func (p *poll) removeFd(fd int32) error {
	if err := syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_DEL, int(fd), nil); err != nil {
		return err
	}

	return nil
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	file, err := ln.(*net.TCPListener).File()
	if err != nil {
		panic(err)
	}

	epoll, err := newPoll(0)
	if err != nil {
		panic(err)
	}

	err = epoll.addEvents(int32(file.Fd()), syscall.EPOLLIN|EPOLLEXCLUSIVE)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(p *poll, n int) {
			accept(p, i)
			wg.Done()
		}(epoll, i)
	}

	wg.Wait()
}

func accept(p *poll, n int) {
	fds, err := p.wait(-1)
	if err != nil {
		fmt.Printf("panic\n")
		panic(err)
	}

	fmt.Println(n, fds)
}