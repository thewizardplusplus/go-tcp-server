package tcpServer

import (
	"net"
)

// is only used to generate the mock
type netConn interface { //nolint:unused
	net.Conn
}

// is only used to generate the mock
type errorHandler interface { //nolint:unused
	HandleError(err error)
}
