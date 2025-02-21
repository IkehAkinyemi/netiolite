package main

import (
	"net"
	"syscall"
)

type conn struct {
	write bool
	data []byte
	oidx int
	poll *poll
	laddr net.Addr
	raddr net.Addr
	saddr syscall.Sockaddr
	fd int
}

func (c *conn) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (c *conn) Close() error {
	return nil
}
