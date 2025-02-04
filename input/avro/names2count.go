package avro

import (
	"context"
	"iter"

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
