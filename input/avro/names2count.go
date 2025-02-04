package avro

import (
	"context"
	"iter"
	"log"
	"runtime"

	ac "github.com/takanoriyanagitani/go-avro-count"
	. "github.com/takanoriyanagitani/go-avro-count/util"
)

func FilenamesToCount(
	filenames iter.Seq[string],
	counter func(string) IO[ac.Count],
) IO[ac.Count] {
	return func(ctx context.Context) (ac.Count, error) {
		var tot ac.Count
		for filename := range filenames {
			cnt, e := counter(filename)(ctx)
			if nil != e {
				return 0, e
			}
			tot += cnt
		}
		return tot, nil
	}
}

func FilenamesToCountSingle(filenames iter.Seq[string]) IO[ac.Count] {
	return FilenamesToCount(
		filenames,
		FilenameToAvroToCount,
	)
}

func FilenamesToCountMulti(
	filenames iter.Seq[string],
	counter func(string) IO[ac.Count],
	workers int,
) IO[ac.Count] {
	return func(ctx context.Context) (ac.Count, error) {
		names := make(chan string)

		go func() {
			defer close(names)

			for name := range filenames {
				names <- name
			}
		}()

		var counts []chan ac.Count
		for range workers {
			counts = append(counts, make(chan ac.Count))
		}

		for i := range workers {
			go func(wi int) {
				var ch chan ac.Count = counts[wi]
				defer close(ch)

				var tot ac.Count
				for filename := range names {
					cnt, e := counter(filename)(ctx)
					if nil != e {
						log.Printf("unable to count: %v\n", e)
						return
					}
					tot += cnt
				}

				ch <- tot
			}(i)
		}

		var tot ac.Count
		for wi := range workers {
			var ch chan ac.Count = counts[wi]
			var cnt ac.Count = <-ch
			tot += cnt
		}

		return tot, nil
	}
}

func FilenamesToCountMultiDefault(filenames iter.Seq[string]) IO[ac.Count] {
	return FilenamesToCountMulti(
		filenames,
		FilenameToAvroToCount,
		runtime.GOMAXPROCS(-1),
	)
}
