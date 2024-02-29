package chmux

import (
	"context"
)

type Mux[T any] struct {
	sinks     []chan T
	ctx       context.Context
	cancel    context.CancelFunc
	c         chan T
	chAddSink chan chan T
	chDelSink chan chan T
}

func (m *Mux[T]) run() {
	for {
		select {
		case <-m.ctx.Done():
			return

		case msg := <-m.chAddSink:
			m.sinks = append(m.sinks, msg)

		case msg := <-m.chDelSink:
			for i, out := range m.sinks {
				if out == msg {
					m.sinks = append(m.sinks[:i], m.sinks[i+1:]...)
					close(msg)
					break
				}
			}

		case msg := <-m.c:
			for _, out := range m.sinks {
				out <- msg
			}
		}
	}
}

func (m *Mux[T]) Close() {
	m.CloseSink(m.sinks...)
	m.cancel()
	close(m.c)
}

func (m *Mux[T]) Send(t T) {
	m.c <- t
}

func (m *Mux[T]) NewSink() chan T {
	ret := make(chan T)
	m.chAddSink <- ret
	return ret
}

func (m *Mux[T]) CloseSink(chs ...chan T) {
	for _, ch := range chs {
		m.chDelSink <- ch
	}
}

func (m *Mux[T]) CloseAllSinks() {
	m.CloseSink(m.sinks...)
}

func New[T any](ctx context.Context) *Mux[T] {
	ctx, cancel := context.WithCancel(ctx)

	ret := &Mux[T]{
		sinks:     nil,
		ctx:       ctx,
		cancel:    cancel,
		c:         make(chan T),
		chAddSink: make(chan chan T),
		chDelSink: make(chan chan T),
	}

	go ret.run()

	return ret
}
