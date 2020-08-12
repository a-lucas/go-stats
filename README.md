## Fun Fact: What is the correct answer ?

https://stats.stackexchange.com/questions/25848/how-to-sum-a-standard-deviation

## Benchmark Results

The native internal stats Golang package is faster than this package.
https://github.com/golang/perf/blob/master/internal/stats/ 

But it doesn't cache results, calling `Sample.Mean()` add extra time to the execution. 


```shell script
BenchmarkSampleStream/Mono_Routine
BenchmarkSampleStream/Mono_Routine/AppendMany_Mono_routine
BenchmarkSampleStream/Mono_Routine/AppendMany_Mono_routine-16         	  398772	      2766 ns/op	       0 B/op	       0 allocs/op
BenchmarkSampleStream/Mono_Routine/Append_Mono_routine
BenchmarkSampleStream/Mono_Routine/Append_Mono_routine-16             	  447692	      2681 ns/op	       0 B/op	       0 allocs/op
BenchmarkSampleStream/Vs_Sample
BenchmarkSampleStream/Vs_Sample-16                                    	  112311	      9974 ns/op	    8191 B/op	      10 allocs/op
BenchmarkSampleStream/Sample_Stream_Multiple_Routine_-_process_1000_
BenchmarkSampleStream/Sample_Stream_Multiple_Routine_-_process_1000_/Concurrency_5
BenchmarkSampleStream/Sample_Stream_Multiple_Routine_-_process_1000_/Concurrency_5-16         	    1082	   1079397 ns/op	     985 B/op	       3 allocs/op
BenchmarkSampleStream/Sample_Stream_Multiple_Routine_-_process_1000_/Concurrency_10
BenchmarkSampleStream/Sample_Stream_Multiple_Routine_-_process_1000_/Concurrency_10-16        	    1146	   1070649 ns/op	     900 B/op	       3 allocs/op
BenchmarkSampleStream/Sample_Stream_Multiple_Routine_-_process_1000_/Concurrency_40
BenchmarkSampleStream/Sample_Stream_Multiple_Routine_-_process_1000_/Concurrency_40-16        	    1143	   1052828 ns/op	     945 B/op	       3 allocs/op
BenchmarkSampleStream/Sample_Stream_Multiple_Routine_-_process_1000_/Concurrency_100
BenchmarkSampleStream/Sample_Stream_Multiple_Routine_-_process_1000_/Concurrency_100-16       	    1170	   1061057 ns/op	     919 B/op	       3 allocs/op
BenchmarkSampleStream/Sample_Multiple_Routine_-_process_1000_
BenchmarkSampleStream/Sample_Multiple_Routine_-_process_1000_/Concurrency_5
BenchmarkSampleStream/Sample_Multiple_Routine_-_process_1000_/Concurrency_5-16                	     346	   3434745 ns/op	 8192633 B/op	   10039 allocs/op
BenchmarkSampleStream/Sample_Multiple_Routine_-_process_1000_/Concurrency_10
BenchmarkSampleStream/Sample_Multiple_Routine_-_process_1000_/Concurrency_10-16               	     346	   3435225 ns/op	 8192561 B/op	   10038 allocs/op
BenchmarkSampleStream/Sample_Multiple_Routine_-_process_1000_/Concurrency_40
BenchmarkSampleStream/Sample_Multiple_Routine_-_process_1000_/Concurrency_40-16               	     385	   3142191 ns/op	 8194301 B/op	   10060 allocs/op
BenchmarkSampleStream/Sample_Multiple_Routine_-_process_1000_/Concurrency_100
BenchmarkSampleStream/Sample_Multiple_Routine_-_process_1000_/Concurrency_100-16              	     379	   3134123 ns/op	 8194583 B/op	   10062 allocs/op

```