package sgo_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/mkch/sgo"
)

func ExampleGroup() {
	sgo.NewGroup().
		Go(func() {
			fmt.Println("task 1")
		}).
		Go(func() {
			fmt.Println("task 2")
		}).
		Wait()
	// Unordered Output:
	// task 1
	// task 2
}

func ExampleWait() {
	sgo.Wait(
		func() { fmt.Println("task 1") },
		func() { fmt.Println("task 2") })

	// Unordered Output:
	// task 1
	// task 2
}

func ExampleGroup_MaxGo() {
	serveConn := func(conn net.Conn) {
		// Serve the conn
	}
	runServer := func(listener net.Listener) {
		// The max concurrency is 2
		group := sgo.NewGroup().MaxGo(2)
		for {
			conn, err := listener.Accept()
			if err != nil {
				break
			}
			group.Go(func() {
				serveConn(conn)
			})
		}
		group.Wait()
	}

	_ = runServer
}

func ExampleGroup_reuse() {
	group := sgo.NewGroup()
	group.MaxGo(2).
		Go(func() { fmt.Println("task 1") }).
		Go(func() { fmt.Println("task 2") }).
		Wait()
	// After group.Wait() returns, group can be reused.
	group.MaxGo(0).
		Go(func() { fmt.Println("task 3") }).
		Go(func() { fmt.Println("task 4") }).
		Wait()

	// The output is unordered, but
	// 1 and 2 are always precede 3 and 4.

	// Unordered Output:
	// task 1
	// task 2
	// task 3
	// task 4
}

func ExampleCollector() {
	sumSlice := func(s []int) (sum int) {
		for _, n := range s {
			sum += n
		}
		return
	}

	s := []int{1, 2, 3, 4}

	result := sgo.NewCollector[int]().
		Go(func() int { return sumSlice(s[:2]) }).
		Go(func() int { return sumSlice(s[2:]) }).
		Collect()

	sum := result[0] + result[1]
	fmt.Println(sum)
	// Output: 10
}

func ExampleCollect() {
	sumSlice := func(s []int) (sum int) {
		for _, n := range s {
			sum += n
		}
		return
	}

	s := []int{1, 2, 3, 4}

	result := sgo.Collect(
		func() int { return sumSlice(s[:2]) },
		func() int { return sumSlice(s[2:]) })
	sum := result[0] + result[1]
	fmt.Println(sum)
	// Output: 10
}

func ExampleRacer() {
	result := sgo.NewRacer[string](context.Background()).
		Go(func(ctx context.Context) string {
			return "1"
		}).
		Go(func(ctx context.Context) string {
			time.Sleep(time.Millisecond * 20)
			// this task fails because of the sleeping
			select {
			case <-ctx.Done():
				return "failed"
			default:
				return "2"
			}
		}).
		Collect()
	fmt.Println(result)
	// Output:
	// 1
}

func ExampleRacer_cleanup() {
	type Result struct {
		Response *http.Response
		Error    error
	}

	racer := sgo.NewRacer[Result](context.Background())
	// There is a chance that more than one racer wins at the same time.
	racer.SetCleanup(func(r Result) {
		if r.Response != nil {
			r.Response.Body.Close()
		}
	})

	for _, url := range []string{"https://go.dev/", "https://golang.org"} {
		racer.Go(func(ctx context.Context) Result {
			request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				return Result{nil, err}
			}
			response, err := http.DefaultClient.Do(request)
			return Result{response, err}
		})
	}

	winner := racer.Collect()
	if winner.Response != nil {
		defer winner.Response.Body.Close()
		// read winner.Response.Body here
	}

}

func ExampleRace() {
	result := sgo.Race(context.Background(),
		func(ctx context.Context) string {
			return "1"
		},
		func(ctx context.Context) string {
			time.Sleep(time.Millisecond * 20)
			// this task fails because of the sleeping
			select {
			case <-ctx.Done():
				return "failed"
			default:
				return "2"
			}
		})
	fmt.Println(result)
	// Output:
	// 1
}
