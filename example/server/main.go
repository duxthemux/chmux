package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/duxthemux/chmux"
)

var mux *chmux.Mux[[]byte]

func handleConn(ctx context.Context, c net.Conn) {
	go func() {
		<-ctx.Done()
		_ = c.Close()
	}()
	go func() {
		out := mux.NewSink()
		defer mux.CloseSink(out)
		for {

			bs, ok := <-out
			if !ok {
				return
			}
			c.Write(bs)
			c.Write([]byte("\n"))
		}
	}()

	go func() {
		scanner := bufio.NewScanner(c)
		for scanner.Scan() {
			txt := scanner.Bytes()
			mux.Send(txt)
			os.Stdout.WriteString(fmt.Sprintf("%s: %s", c.RemoteAddr().String()))
		}
	}()

}
func run(ctx context.Context) error {
	mux = chmux.New[[]byte](ctx)

	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		_ = l.Close()
	}()

	for {
		if ctx.Err() != nil {
			return err
		}
		c, err := l.Accept()
		if err != nil {
			slog.Warn("did not accept", "err", err)
		}

		go handleConn(ctx, c)
	}

	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		panic(err)
	}
}
