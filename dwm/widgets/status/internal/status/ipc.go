package status

import (
	"context"
	"errors"
	"net"
	"os"
	"strings"
	"time"
)

func Serve(ctx context.Context, socket string, d *Daemon) error {
	_ = os.Remove(socket)
	addr, err := net.ResolveUnixAddr("unixgram", socket)
	if err != nil {
		return err
	}
	conn, err := net.ListenUnixgram("unixgram", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	defer os.Remove(socket)

	buf := make([]byte, 128)
	for {
		_ = conn.SetReadDeadline(time.Now().Add(time.Second))
		n, _, err := conn.ReadFromUnix(buf)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					continue
				}
			}
			return err
		}
		msg := strings.TrimSpace(string(buf[:n]))
		parts := strings.Fields(msg)
		if len(parts) == 2 && parts[0] == "refresh" {
			_ = d.Refresh(ctx, parts[1])
		} else if msg == "refresh" || msg == "" {
			_ = d.Refresh(ctx, "all")
		}
	}
}

func Send(socket string, msg string, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = time.Second
	}
	addr, err := net.ResolveUnixAddr("unixgram", socket)
	if err != nil {
		return err
	}
	conn, err := net.DialUnix("unixgram", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	_ = conn.SetWriteDeadline(time.Now().Add(timeout))
	_, err = conn.Write([]byte(msg))
	return err
}

func Health(socket string, timeout time.Duration) error {
	if _, err := os.Stat(socket); err != nil {
		return err
	}
	if err := Send(socket, "health", timeout); err != nil {
		return err
	}
	return nil
}

var ErrNoDaemon = errors.New("dwm-status daemon is not available")
