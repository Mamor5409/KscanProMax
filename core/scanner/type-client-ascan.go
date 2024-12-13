package scanner

import (
	"KscanPro/core/hydra"
	"net"
)

type AsacnClinet struct {
	*client
	// 后面改下
	HandlerSuccess func(addr net.IP, port int, protocol string, auth *hydra.Auth)
	HandlerError   func(domain string, err error)
}

func NewAscancanner(config *Config) *DomainClient {

	return nil
}
