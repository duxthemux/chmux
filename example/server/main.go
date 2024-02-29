package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"

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
			_, err := c.Write(bs)
			if err != nil {
				slog.Warn("could not write to con", "err", err)
				return
			}
			_, err = c.Write([]byte("\n"))
			if err != nil {
				slog.Warn("could not write to con", "err", err)
				return
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(c)
		for scanner.Scan() {
			txt := scanner.Text()
			pl := fmt.Sprintf("%s: %s", c.RemoteAddr().String(), txt)
			mux.Send([]byte(pl))
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
}

func main() {
	if err := run(context.Background()); err != nil {
		panic(err)
	}
}
