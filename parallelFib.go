package main

import ("fmt"
"sync")

func main() {
    in := gen(6, 10)

    // Distribute the sq work across two goroutines that both read from in.
    c1 := makeServer(in)
    c2 := makeServer(in)


	
    for n := range merge(c1, c2) {
        fmt.Println(n) // 4 then 9, or 9 then 4
    }
}

func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func merge(cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    // Start an output goroutine for each input channel in cs.  output
    // copies values from c to out until c is closed, then calls wg.Done.
    output := func(c <-chan int) {
        for n := range c {
            out <- n
        }
        wg.Done()
    }
    wg.Add(len(cs))
    for _, c := range cs {
        go output(c)
    }

    // Start a goroutine to close out once all the output goroutines are
    // done.  This must start after the wg.Add call.
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}

func makeServer(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- fib(n)
	}
        close(out)
    }()
    return out
}

func fib(a int) int {
	if a < 2 {
		return a
	}
	return fib(a - 1) + fib(a - 2)
}
