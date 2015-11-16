package main

import (
  "fmt"
  "sync"
  "time"
  "math/rand"
)

func main() {
    //in := makeCustomers(39,45,46,35)
    //in := makeCustomers(29,35,36,25)

    numCustomer := 4
    in := makeRandomCustomers(numCustomer)

    // Distribute the sq work across two goroutines that both read from in.
    c1 := makeServer(in)
    c2 := makeServer(in)
  //  c3 := makeServer(in)



    for n := range merge(c1 ,c2) {
        fmt.Println(n)
    }

}

func makeRandomCustomers(count int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    wg.Add(count);//wait for all count go procedures to finish

    for i := 0; i < count; i++ {
      go func(index int) {
        fmt.Printf("Starting customer %d\n",index)
        defer fmt.Printf("Customer %d is done!\n",index)
        defer wg.Done() //decrement wg when this finishes
        //fmt.Printf("start %d\n",index)
        delaySeed := rand.Intn(20)
        delay := time.Duration(delaySeed)*time.Second
        time.Sleep(delay)
        out <- rand.Intn(20)+10 //Passing a random number to get fib-ed
      }(i)//pass in a copy of index so we can keep track of which customer
    }

  //  wg.Wait()
    //close(out)
    return out
}


func makeCustomers(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums { //not using the key, only the value
                  fmt.Println("start")
                  duration := time.Duration(2)*time.Second // Pause for 10 seconds
                  time.Sleep(duration)
                  fmt.Println("done")
                  fmt.Println(n)
            out <- n
        }
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


func fib(a int) int {
	if a < 2 {
		return a
	}
	return fib(a - 1) + fib(a - 2)
}
