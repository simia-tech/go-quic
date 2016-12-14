package quic_test

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func ListenUDP(tb testing.TB, localAddress string) *net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp", localAddress)
	require.NoError(tb, err)
	conn, err := net.ListenUDP("udp", addr)
	require.NoError(tb, err)
	return conn
}

func DialUDP(tb testing.TB, addr net.Addr) *net.UDPConn {
	udpAddr, ok := addr.(*net.UDPAddr)
	require.True(tb, ok)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	require.NoError(tb, err)
	return conn
}

func ReadLine(tb testing.TB, r io.Reader) string {
	line, err := bufio.NewReader(r).ReadString('\n')
	require.NoError(tb, err)
	return line
}

func WriteLine(tb testing.TB, w io.Writer, line string) {
	_, err := fmt.Fprintf(w, "%s\n", line)
	require.NoError(tb, err)
}
