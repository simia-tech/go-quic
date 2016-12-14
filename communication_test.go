package quic_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/go-quic"
)

func TestClientServerEcho(t *testing.T) {
	serverConn := ListenUDP(t, "localhost:0")

	listener, err := quic.Listen(serverConn, 3)
	require.NoError(t, err)

	go func() {
		conn, err := listener.Accept()
		require.NoError(t, err)
		defer conn.Close()

		line := ReadLine(t, conn)
		WriteLine(t, conn, line)
	}()

	clientConn := DialUDP(t, serverConn.LocalAddr())

	conn, err := quic.Dial(clientConn, 3)
	require.NoError(t, err)

	WriteLine(t, conn, "test")
	line := ReadLine(t, conn)

	assert.Equal(t, "test", line)
}
