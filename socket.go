package main

import (
	"os"
	"syscall"
)

type listener struct {
	fd int
	sa *syscall.SockaddrInet4
}

func NewListener(ip []byte, port int, flags int) (*listener, error) {
	syscall.ForkLock.Lock()
	
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM | flags, syscall.IPPROTO_TCP)
	if err != nil {
		return nil, os.NewSyscallError("socket", err)
	}

	syscall.ForkLock.Unlock()

	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		syscall.Close(fd)
		return nil, os.NewSyscallError("setsockoptInt", err)
	}

	sa := &syscall.SockaddrInet4{Port: port}
	copy(sa.Addr[:], ip)

	if err := syscall.Bind(fd, sa); err != nil {
		syscall.Close(fd)
		return nil, os.NewSyscallError("bind", err)
	}

	if err := syscall.Listen(fd, syscall.SOMAXCONN); err != nil {
		syscall.Close(fd)
		return nil, os.NewSyscallError("listen", err)
	}

	return &listener{fd, sa}, nil
}

func (so *listener) Accept(flags int) (*conn, error) {
	nfd, sa, err := syscall.Accept4(so.fd, flags)
	if err != nil && err != syscall.EAGAIN {
		return nil, os.NewSyscallError("accept", err)
	}

	return &conn{fd: nfd, saddr: sa}, nil
}