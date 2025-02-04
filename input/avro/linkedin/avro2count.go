package avro2count

import (
	"bufio"
	"context"
	"io"
	"os"

	lg "github.com/linkedin/goavro/v2"
	ac "github.com/takanoriyanagitani/go-avro-count"
	util "github.com/takanoriyanagitani/go-avro-count/util"
)

func CountAvro(ctx context.Context, r io.Reader) (ac.Count, error) {
	var br io.Reader = bufio.NewReader(r)

	rdr, e := lg.NewOCFReader(br)
	if nil != e {
		return 0, e
	}

	var cnt uint64 = 0
	for rdr.Scan() {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}

		cnt += uint64(rdr.RemainingBlockItems())

		rdr.SkipThisBlockAndReset()
	}

	return ac.Count(cnt), nil
}

func ReaderToCount(r io.Reader) util.IO[ac.Count] {
	return func(ctx context.Context) (ac.Count, error) {
		return CountAvro(ctx, r)
	}
}

var StdinToAvroToCount util.IO[ac.Count] = ReaderToCount(os.Stdin)
