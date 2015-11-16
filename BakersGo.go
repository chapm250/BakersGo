package main

import (
  "fmt"
  "sync"
  "time"
  "math/rand"
)

type customer struct {
	id int
	fib int
}

func main() {
    numCustomers := 10;



    in := makeRandomCustomers(numCustomers)

    // Distribute the sq work across two goroutines that both read from in.
    // c1 := serialServer(in)
    // c2 := serialServer(in)
  //  c3 := makeServer(in)

// fmt.Println(c1)
// fmt.Println(c2)

  allResults := merge(serialServer(in),serialServer(in))
  for i := 0; i < numCustomers; i++ {
    doneCust := <-allResults

    fmt.Printf("========>The result is %d for customer %d\n", doneCust.fib, doneCust.id)
  }
    // terminator.Done()
    // terminator.Wait()
}

func makeRandomCustomers(count int) <-chan customer {
    out := make(chan customer)
  //  var wg sync.WaitGroup
//    wg.Add(count);//wait for all count go procedures to finish

    for i := 0; i < count; i++ {
      go func(index int) {
        fmt.Printf("Starting customer %d\n",index)
    //    defer fmt.Printf("Customer %d is done!\n",index)
    //    defer wg.Done() //decrement wg when this finishes
        delay := time.Duration(rand.Intn(20))*time.Second
        fibNum := rand.Intn(40)+10;
        time.Sleep(delay)
        fmt.Printf("Done sleeping. Calculate fib(%d) for %d\n",fibNum, index)
        out <- customer{index, fibNum} //Passing a random number to get fib-ed
      }(i)//pass in a copy of index so we can keep track of which customer
    }

  //  wg.Wait()
    //close(out)
    return out
}


func serialServer(in <-chan customer) <-chan customer {
    out := make(chan customer)
    go func() {
        for n := range in {
            fmt.Printf("-->Calculating fib of %d for customer %d\n", n.fib, n.id)
            out <- customer{n.id, fib(n.fib)}
            fmt.Printf("?=====>Done Calculating fib of %d for customer %d??\n", n.fib, n.id)
	}
        close(out)
    }()
    return out
}


func merge(cs ...<-chan customer) <-chan customer {
    var wg sync.WaitGroup
    out := make(chan customer)

    // Start an output goroutine for each input channel in cs.  output
    // copies values from c to out until c is closed, then calls wg.Done.
    output := func(c <-chan customer) {
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
