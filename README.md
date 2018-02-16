# GoDNA DNA generator

GoDNA is a storage benchmarking/load-generation tool that generates a random DNA sequence and writes it concurently to multiple files.  

## Usage

```
Usage: ./godna [-t <cpu threads> -c <concurrency> -d </dir/path>] -n <number of files> -s <size in bytes>

  -n int
    	
      Required: set the number of files to generate.
  
  -s int

    	Required: the size of the file(s) to generate in Bytes (1MiB = 1048576 Bytes)

  -c int
    	
      Optional: set number of concurrent file writers to use; defaults to the number of logical CPUs in your system.
      Setting this allows you to specify how many files will be writen at the same time. By default it will write as
      many files concurrently as there are logical CPUs in your system. 
  
  -d string

    	Optional: Directory to write the files in; defaults to the current directory (default ".")
  
  -t int

      Optional: set the number CPU threads available for use; defaults to the number of logical CPUs in your system.
      This setting lets you specify how many CPU threads will used to schedule the concurrent IO workers on. 
  
  -v	

      Optional: Turn on verbose output mode; it will print the progress every second.
```

## Examples

### Example 1

Create 10000 1MiB files in the "/data/dna" directory using as many threads and writers as logical CPUs in your system: 

```
$ ./godna -n 10000 -s 1048576 -d /data/dna

Done!
Number of Files Written: 10000, Total Size: 9.8GiB, Avg FPS: 333, Avg Throughput: 333 MiB/s, Elapsed Time: 30 seconds
```

### Example 2

Create 1000 10MiB files in the "/data/dna" directory using 16 CPU threads with 32 concurent IO writers (two writers per thread) with verbose output:

```
$ ./godna -v -n 1024 -t 16 -c 32 -s 10485760 -d /data/dna

Files Completed:      85, Data Written:   0.8GiB, Files Remaining:     949, Cur FPS:    85, Throughput:  850 MiB/s
Files Completed:     179, Data Written:   1.7GiB, Files Remaining:     851, Cur FPS:    89, Throughput:  895 MiB/s
Files Completed:     277, Data Written:   2.7GiB, Files Remaining:     759, Cur FPS:    92, Throughput:  923 MiB/s
Files Completed:     360, Data Written:   3.5GiB, Files Remaining:     670, Cur FPS:    90, Throughput:  900 MiB/s
Files Completed:     467, Data Written:   4.6GiB, Files Remaining:     565, Cur FPS:    93, Throughput:  934 MiB/s
Files Completed:     564, Data Written:   5.5GiB, Files Remaining:     466, Cur FPS:    94, Throughput:  940 MiB/s
Files Completed:     662, Data Written:   6.5GiB, Files Remaining:     374, Cur FPS:    94, Throughput:  945 MiB/s
Files Completed:     747, Data Written:   7.3GiB, Files Remaining:     283, Cur FPS:    93, Throughput:  933 MiB/s
Files Completed:     848, Data Written:   8.3GiB, Files Remaining:     185, Cur FPS:    94, Throughput:  942 MiB/s
Files Completed:     957, Data Written:   9.3GiB, Files Remaining:      76, Cur FPS:    95, Throughput:  957 MiB/s

Done!
Number of Files Written: 1024, Total Size: 10.0GiB, Avg FPS: 93, Avg Throughput: 930 MiB/s, Elapsed Time: 11 seconds
```