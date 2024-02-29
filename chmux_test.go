package chmux_test

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/duxthemux/chmux"
)

func TestSimple(t *testing.T) {
	ctx := context.Background()

	mux := chmux.New[int](ctx)
	defer mux.Close()

	out1 := mux.NewSink()

	total := 0
	count := 0

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		for {
			out, ok := <-out1
			if !ok {
				return
			}

			log.Printf("Out1: %v", out)

			total = total + out
			count++

			wg.Done()

		}
	}()

	go func() {
		for i := 0; i < 4; i++ {
			mux.Send(i)
		}
	}()

	wg.Wait()

	log.Printf("Total: %v", total)
	log.Printf("Count: %v", count)

}

func TestSimple2(t *testing.T) {
	ctx := context.Background()
	mux := chmux.New[int](ctx)
	defer mux.Close()

	out1 := mux.NewSink()
	out2 := mux.NewSink()

	total := 0
	count := 0

	wg := sync.WaitGroup{}
	wg.Add(8)
	go func() {
		for {
			out, ok := <-out1
			if !ok {
				return
			}

			log.Printf("Out1: %v", out)

			total = total + out
			count++
			wg.Done()

		}
	}()

	go func() {
		for {
			out, ok := <-out2
			if !ok {
				return
			}

			log.Printf("Out2: %v", out)

			total = total + out
			count++
			wg.Done()

		}
	}()

	go func() {
		for i := 0; i < 4; i++ {
			mux.Send(i)
		}
	}()

	wg.Wait()

	log.Printf("Total: %v", total)
	log.Printf("Count: %v", count)

}

func Test2In(t *testing.T) {
	ctx := context.Background()

	mux := chmux.New[int](ctx)
	defer mux.Close()

	out1 := mux.NewSink()

	total := 0
	count := 0

	wg := sync.WaitGroup{}
	wg.Add(8)
	go func() {
		for {
			out, ok := <-out1
			if !ok {
				return
			}

			log.Printf("Out1: %v", out)

			total = total + out
			count++
			wg.Done()

		}
	}()

	go func() {
		for i := 0; i < 4; i++ {
			mux.Send(i)
		}

	}()

	go func() {
		for i := 0; i < 4; i++ {
			mux.Send(i)
		}

	}()

	wg.Wait()

	log.Printf("Total: %v", total)
	log.Printf("Count: %v", count)

}

func Test2In2Out(t *testing.T) {
	ctx := context.Background()

	mux := chmux.New[int](ctx)
	defer mux.Close()

	out1 := mux.NewSink()
	out2 := mux.NewSink()

	total := 0
	count := 0
	writes := 0

	wg := sync.WaitGroup{}

	wg.Add(16)
	go func() {
		for {
			out, ok := <-out1
			if !ok {
				return
			}

			log.Printf("Out1: %v", out)

			total = total + out
			count++
			time.Sleep(time.Millisecond * time.Duration(count))
			wg.Done()

		}
	}()

	go func() {
		for {
			out, ok := <-out2
			if !ok {
				return
			}

			log.Printf("Out2: %v", out)

			total = total + out
			count++
			time.Sleep(time.Millisecond * time.Duration(count))
			wg.Done()

		}
	}()

	go func() {
		for i := 0; i < 4; i++ {
			mux.Send(10 + i)
			writes++
		}

	}()

	go func() {
		for i := 0; i < 4; i++ {
			mux.Send(20 + i)
			writes++
		}

	}()

	wg.Wait()
	mux.CloseSink(out1, out2)

	log.Printf("Total: %v", total)
	log.Printf("Count: %v", count)
	log.Printf("Writes: %v", writes)

}

func TestStepByStep(t *testing.T) {
	ctx := context.Background()

	wg := sync.WaitGroup{}
	mux := chmux.New[int](ctx)
	defer mux.Close()

	t.Run("Step1", func(t *testing.T) {

		wg.Add(1)

		out1 := mux.NewSink()

		go func() {
			i := <-out1
			log.Printf("Got: %v", i)
			wg.Done()
		}()

		mux.Send(1)
		wg.Wait()
		mux.CloseSink(out1)
	})

	t.Run("Problematic", func(t *testing.T) {

		wg = sync.WaitGroup{}
		wg.Add(2)

		out1 := mux.NewSink()
		out2 := mux.NewSink()

		go func() {
			for {
				i, ok := <-out1
				if !ok {
					return
				}
				log.Printf("Got: %v", i)
				wg.Done()
			}
		}()

		go func() {
			for {
				i, ok := <-out2
				if !ok {
					return
				}
				log.Printf("Got: %v", i)
				wg.Done()
			}
		}()

		mux.Send(2)

		wg.Wait()
		mux.CloseAllSinks()
	})
}
