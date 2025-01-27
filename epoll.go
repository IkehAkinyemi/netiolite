package main

import (
	"syscall"
	"time"
)

type poll struct {
	fd int
	events []syscall.EpollEvent
	evfds []int32
}

func newPoll() *poll {
	fd, err := syscall.EpollCreate1(0)
	if err != nil {
		panic(err)
	}

	p := new(poll)
	p.fd = fd

	p.events = make([]syscall.EpollEvent, 0, 64)
	p.evfds = make([]int32, len(p.events))

	return p
}

func (p *poll) wait(msec time.Duration) []int32 {
	var err error
	var n int
	if msec > 0 {
		n, err = syscall.EpollWait(p.fd, p.events, int(msec * time.Millisecond))
	} else if msec == 0 {
		n, err = syscall.EpollWait(p.fd, p.events, 0)
	} else {
		n, err = syscall.EpollWait(p.fd, p.events, -1)
	}

	if err != nil && err != syscall.EINTR {
		panic(err)
	}

	p.evfds = p.evfds[:0]
	for i := 0; i < n; i++ {
		p.evfds = append(p.evfds, p.events[i].Fd)
	}

	return p.evfds
}

func (p *poll) addRead(fd int32) {
	if err := syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_ADD, int(fd), &syscall.EpollEvent{
		Fd: fd,
		Events: syscall.EPOLLIN,
	}); err != nil {
		panic(err)
	}
}

func (p *poll) modReadWrite(fd int32) {
	if err := syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_MOD, int(fd), &syscall.EpollEvent{
		Fd: fd,
		Events: syscall.EPOLLIN | syscall.EPOLLOUT,
	}); err != nil {
		panic(err)
	}
}

func (p *poll) modRead(fd int32) {
	if err := syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_MOD, int(fd), &syscall.EpollEvent{
		Fd: fd,
		Events: syscall.EPOLLIN,
	}); err != nil {
		panic(err)
	}
}