package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {

	// define and set default command parameter flags
	var nFlag = flag.Int("n", 0, "Required: set the number of files to generate")
	var sFlag = flag.Int("s", 0, "Required: the size of the file(s) to generate in Bytes")
	var tFlag = flag.Int("t", runtime.NumCPU(), "Optional: set number CPU threads available for use; defaults to the number of logical CPUs in your system")
	var cFlag = flag.Int("c", runtime.NumCPU(), "Optional: set number of concurrent file writers to use; defaults to the number of logical CPUs in your system")
	var dFlag = flag.String("d", ".", "Optional: Directory to write the files; defaults to the current directory")
	var vFlag = flag.Bool("v", false, "Optional: Turn on verbose output mode; it will print the progres every second")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage: %s [-t <cpu threads> -c <concurrency> -d </dir/path>] -n <number of files> -s <size in bytes>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExample: %s -n 1024 -t 16 -c 32 -s 10485760 -d /tmp\n\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println()
	}
	flag.Parse()

	if *sFlag == 0 {
		fmt.Fprintf(os.Stderr, "\nMissing required -s (size) argument\n\n")
		flag.Usage()
		os.Exit(2)
	}

	if *nFlag == 0 {
		fmt.Fprintf(os.Stderr, "\nMissing required -n (number of files) argument\n\n")
		flag.Usage()
		os.Exit(2)
	}

	runtime.GOMAXPROCS(*tFlag)
	size := *sFlag
	wg := new(sync.WaitGroup)
	sema := make(chan struct{}, *cFlag)
	out := genstring(size)
	progress := make(chan int64, 256)
	start := time.Now().Unix()


        go func() {
                wg.Wait()
                close(progress)
        }()

	for x := 1; x <= *nFlag; x++ {
		wg.Add(1)
		go spraydna(x, wg, sema, &out, *dFlag, progress)
	}

        // If the '-v' flag was provided, periodically print the progress stats
        var tick <-chan time.Time
        if *vFlag {

	        tick = time.Tick(1000 * time.Millisecond)
        }

        var nfiles, nbytes int64
	
loop:
        for {
                select {
                case size, ok := <-progress:
                        if !ok {
                        	break loop // progress was closed
                        }
                        nfiles++
                        nbytes += size
                case <-tick:
                        printProgress(nfiles, nbytes, start)
                }
        }

        // Final totals
        printDiskUsage(nfiles, nbytes, start)

}


func genstring(size int) []byte {
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Seed(time.Now().UnixNano())
	dnachars := []byte("GATC")
	dna := make([]byte, 0)
	for x := 0; x < size; x++ {
		dna = append(dna, dnachars[rand.Intn(len(dnachars))])
	}
	return dna
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func spraydna(count int, wg *sync.WaitGroup, sema chan struct{}, out *[]byte, dir string, progress chan<- int64) {
	defer wg.Done()
	sema <- struct{}{}        // acquire token
	defer func() { <-sema }() // release token
	filename := fmt.Sprintf("%s/dna-%d.txt", dir, count)
	f, err := os.Create(filename)
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	writtenBytes, _ := w.Write(*out)
	//fmt.Printf("wrote %d bytes\n", writtenBytes)
	w.Flush()
	progress <- int64(writtenBytes) 
}

func printProgress(nfiles, nbytes int64, start int64) {
        now := time.Now().Unix()
        elapsed := now - start
        if elapsed == 0 {
                elapsed = 1
        }
        fps := nfiles / elapsed
        tp := nbytes / elapsed
        fmt.Printf("Files Completed: %7d, Data Written: %5.1fGiB, Files Remaining: %7d, Cur FPS: %5d, Throughput: %4d MiB/s\n", nfiles, float64(nbytes)/1073741824, runtime.NumGoroutine(), fps, tp/1048576)
}

// Prints the final summary
func printDiskUsage(nfiles, nbytes int64, start int64) {
        stop := time.Now().Unix()
        elapsed := stop - start
        if elapsed == 0 {
                elapsed = 1
        }
        fps := nfiles / elapsed
        tp := nbytes / elapsed
        fmt.Printf("\nDone!\nNumber of Files Written: %d, Total Size: %.1fGiB, Avg FPS: %d, Avg Throughput: %d MiB/s, Elapsed Time: %d seconds\n", nfiles, float64(nbytes)/1073741824, fps, tp/1048576, elapsed)
}
