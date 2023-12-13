package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
    "time"
)

type RingBuffer struct {
    data     []int
    head     int
    tail     int
    size     int
    interval time.Duration
}

func NewRingBuffer(size int, interval time.Duration) *RingBuffer {
    return &RingBuffer{
        data:     make([]int, size),
        head:     0,
        tail:     0,
        size:     size,
        interval: interval,
    }
}

func (r *RingBuffer) Push(value int) {
    r.data[r.head] = value
    r.head = (r.head + 1) % r.size
    if r.head == r.tail {
        r.tail = (r.tail + 1) % r.size
    }
}

func (r *RingBuffer) Pop() int {
    if r.head == r.tail {
        return -1
    }
    value := r.data[r.tail]
    r.tail = (r.tail + 1) % r.size
    return value
}

func main() {

    negativeFilter := func(in <-chan int) <-chan int {
        out := make(chan int)
        go func() {
            for num := range in {
                if num >= 0 {
                    out <- num
                }
            }
            close(out)
        }()
        return out
    }

    multipleOfThreeFilter := func(in <-chan int) <-chan int {
        out := make(chan int)
        go func() {
            for num := range in {
                if num%3 == 0 && num != 0 {
                    out <- num
                }
            }
            close(out)
        }()
        return out
    }

    
    bufferSize := 5
    flushInterval := time.Second * 2
    buffer := NewRingBuffer(bufferSize, flushInterval)
    bufferStage := func(in <-chan int) <-chan int {
        out := make(chan int)
        go func() {
            for num := range in {
                buffer.Push(num)

                if buffer.Pop() != -1 {
                    // If buffer is not empty, flush the data
                    fmt.Println("Flushing data:", buffer.data)
                    buffer.Reset()
                }
            }
            close(out)
        }()
        return out
    }

    source := func() <-chan int {
        out := make(chan int)
        scanner := bufio.NewScanner(os.Stdin)
        go func() {
            defer close(out)
            for scanner.Scan() {
                text := scanner.Text()
                num, err := strconv.Atoi(text)
                if err == nil {
                    out <- num
                } else {
                    fmt.Println("Invalid input:", text)
                }
            }
            if err := scanner.Err(); err != nil {
                fmt.Println("Error reading input:", err)
            }
        }()
        return out
    }

    consumer := func(in <-chan int) {
        for num := range in {
            fmt.Println("Received data:", num)
        }
    }

  
    pipeline := negativeFilter(multipleOfThreeFilter(bufferStage(source())))

    
    consumer(pipeline)
}
