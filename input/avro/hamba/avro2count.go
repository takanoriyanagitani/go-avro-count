package avro2count

import (
	"bufio"
	"context"
	"io"
	"os"

	ho "github.com/hamba/avro/v2/ocf"
	ac "github.com/takanoriyanagitani/go-avro-count"
	util "github.com/takanoriyanagitani/go-avro-count/util"
)

func CountAvro(ctx context.Context, r io.Reader) (ac.Count, error) {
	var br io.Reader = bufio.NewReader(r)

	dec, e := ho.NewDecoder(br)
	if nil != e {
		return 0, e
	}

	var cnt uint64 = 0
	var buf map[string]any
	var err error
	for dec.HasNext() {
		clear(buf)
		err = dec.Decode(&buf)
		if nil != err {
			return 0, err
		}

		cnt += 1
	}

	return ac.Count(cnt), nil
}

func ReaderToCount(r io.Reader) util.IO[ac.Count] {
	return func(ctx context.Context) (ac.Count, error) {
		return CountAvro(ctx, r)
	}
}

var StdinToAvroToCount util.IO[ac.Count] = ReaderToCount(os.Stdin)
