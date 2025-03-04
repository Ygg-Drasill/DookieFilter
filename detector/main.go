package main

import (
    "fmt"
    "github.com/Ygg-Drasill/DookieFilter/detector/detec"
    "log"
)

func main() {
    f, err := detec.LoadFrames()
    if err != nil {
        log.Fatal(err)
    }
    d := detec.Detector{
        Swaps:     []detec.Swap{},
        PlayerMap: make(map[int][]float64),
    }
    fmt.Println(d)
    s := d.Detect(f)
    for _, swap := range s {
        log.Printf("Swap between %s and %s at frame %d", swap.P1.Number, swap.P2.Number, swap.SwapFrame)
    }
    fmt.Println("Done")
}
