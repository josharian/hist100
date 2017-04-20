package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	flagPrefix = flag.String("p", "", "prefix to seek and strip")
	flagEvery  = flag.Int("f", 0, "print running hist every f lines")
)

type counts [10]int

func (c counts) print() {
	tot := 0
	for _, n := range c {
		tot += n
	}
	var pct [10]float64
	for i, n := range c {
		pct[i] = float64(n) / float64(tot)
	}
	fmt.Println()
	for x := 9; x >= 0; x-- {
		fmt.Printf("% 3d%%|", x*10)
		for _, n := range pct {
			if n > float64(x)/10 {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	fmt.Println("----+----------")
	fmt.Println("    |0123456789")
	fmt.Println()
	for i, f := range pct {
		fmt.Printf("%d: %0.2f", i, f)
		if i != len(pct)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println()
}

func main() {
	flag.Parse()
	s := bufio.NewScanner(os.Stdin)
	n := 0
	var count counts
	prefix := []byte(*flagPrefix)
	for s.Scan() {
		buf := s.Bytes()
		if !bytes.HasPrefix(buf, prefix) {
			continue
		}
		buf = bytes.TrimPrefix(buf, prefix)
		buf = bytes.TrimSpace(buf)
		s := string(buf)
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "malformed input %q: %v\n", s, err)
			continue
		}
		if f < 0 || f > 1 {
			fmt.Fprintf(os.Stderr, "input must be in range [0, 1], got %v\n", f)
			continue
		}
		if f == 1 {
			f = 0.99
		}
		f *= 10
		count[int(f)]++
		n++
		if *flagEvery > 0 && n%*flagEvery == 0 {
			count.print()
		}
	}
	if err := s.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "stopped at line %d: %v\n", n, err)
	}
	count.print()
}
