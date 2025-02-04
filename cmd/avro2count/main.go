package main

import (
	"bufio"
	"context"
	"fmt"
	"iter"
	"log"
	"os"
	"strconv"
	"strings"

	ac "github.com/takanoriyanagitani/go-avro-count"
	ia "github.com/takanoriyanagitani/go-avro-count/input/avro"
	util "github.com/takanoriyanagitani/go-avro-count/util"
)

func envVarByKey(key string) util.IO[string] {
	return func(_ context.Context) (string, error) {
		val, found := os.LookupEnv(key)
		switch found {
		case true:
			return val, nil
		default:
			return "", fmt.Errorf("env var %s missing", key)
		}
	}
}

var stdinAsFilenames util.IO[bool] = util.Bind(
	envVarByKey("ENV_STDIN_AS_FILENAMES"),
	util.Lift(strconv.ParseBool),
).Or(util.Of(false))

var useMulti util.IO[bool] = util.Bind(
	envVarByKey("ENV_MULTI_THREAD"),
	util.Lift(strconv.ParseBool),
).Or(util.Of(true))

var filenames util.IO[iter.Seq[string]] = func(
	_ context.Context,
) (iter.Seq[string], error) {
	return func(yield func(string) bool) {
		var s *bufio.Scanner = bufio.NewScanner(os.Stdin)
		for s.Scan() {
			var filename string = s.Text()
			if !yield(filename) {
				return
			}
		}
	}, nil
}

var stdin2names2count util.IO[ac.Count] = util.Bind(
	filenames,
	func(names iter.Seq[string]) util.IO[ac.Count] {
		return util.Bind(
			useMulti,
			func(multi bool) util.IO[ac.Count] {
				switch multi {
				case true:
					return ia.FilenamesToCountMultiDefault(names)
				default:
					return ia.FilenamesToCountSingle(names)
				}
			},
		)
	},
)

var stdin2avro2count util.IO[ac.Count] = ia.StdinToAvroToCount

var stdin2count util.IO[ac.Count] = util.Bind(
	stdinAsFilenames,
	func(isFilenames bool) util.IO[ac.Count] {
		switch isFilenames {
		case true:
			return stdin2names2count
		default:
			return stdin2avro2count
		}
	},
)

func CountToStdout(c ac.Count) util.IO[util.Void] {
	return func(_ context.Context) (util.Void, error) {
		fmt.Printf("%v\n", c)
		return util.Empty, nil
	}
}

var stdin2count2stdout util.IO[util.Void] = util.Bind(
	stdin2count,
	CountToStdout,
)

func main() {
	_, e := stdin2count2stdout(context.Background())
	if nil != e {
		var s string = e.Error()
		if strings.Contains(s, "cannot create OCFReader: cannot read OCF header with invalid magic bytes:") {
			log.Printf("hint: you need to set ENV_STDIN_AS_FILENAMES to use stdin as filenames\n")
		}
		log.Printf("%v\n", e)
	}
}
