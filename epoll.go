package netiolite

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

	p.events = make([]syscall.EpollEvent, 64)
	p.evfds = make([]int32, len(p.events))

	return p
}

func (p *poll) wait(msec time.Duration) []int32 {
	var err error
	var n int
	if msec >= 0 {
		n, err = syscall.EpollWait(p.fd, p.events, int(msec / time.Millisecond))
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

func (p *poll) addEvents(fd int32, flags int) {
	if err := syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_ADD, int(fd), &syscall.EpollEvent{
		Fd: fd,
		Events: uint32(flags),
	}); err != nil {
		panic(err)
	}
}

func (p *poll) modEvents(fd int32, flags int) {
	if err := syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_MOD, int(fd), &syscall.EpollEvent{
		Fd: fd,
		Events: uint32(flags),
	}); err != nil {
		panic(err)
	}
}

func (p *poll) removeFd(fd int32) {
	if err := syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_DEL, int(fd), nil); err != nil {
		panic(err)
	}
}