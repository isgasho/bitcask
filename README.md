# bitcask

[![Build Status](https://cloud.drone.io/api/badges/prologic/bitcask/status.svg)](https://cloud.drone.io/prologic/bitcask)
[![CodeCov](https://codecov.io/gh/prologic/bitcask/branch/master/graph/badge.svg)](https://codecov.io/gh/prologic/bitcask)
[![Go Report Card](https://goreportcard.com/badge/prologic/bitcask)](https://goreportcard.com/report/prologic/bitcask)
[![GoDoc](https://godoc.org/github.com/prologic/bitcask?status.svg)](https://godoc.org/github.com/prologic/bitcask) 
[![Sourcegraph](https://sourcegraph.com/github.com/prologic/msgbus/-/badge.svg)](https://sourcegraph.com/github.com/prologic/msgbus?badge)

A Bitcask (LSM+WAL) Key/Value Store written in Go.

## Features

* Embeddable
* Builtin CLI

## Install

```#!bash
$ go get github.com/prologic/bitcask
```

## Usage (library)

Install the package into your project:

```#!bash
$ go get github.com/prologic/bitcask
```

```#!go
package main

import (
    "log"

    "github.com/prologic/bitcask"
)

func main() {
    db, _ := bitcask.Open("/tmp/db")
    db.Set("Hello", []byte("World"))
    db.Close()
}
```

See the [godoc](https://godoc.org/github.com/prologic/bitcask) for further
documentation and other examples.

## Usage (tool)

```#!bash
$ bitcask -p /tmp/db set Hello World
$ bitcask -p /tmp/db get Hello
World
```

## Performance

Benchmarks run on a 11" Macbook with a 1.4Ghz Intel Core i7:

```
$ make bench
...
BenchmarkGet-4   	   50000	     33185 ns/op	     600 B/op	      14 allocs/op
BenchmarkPut-4   	  100000	     16757 ns/op	     699 B/op	       7 allocs/op
```

* ~30,000 reads/sec
* ~60,000 writes/sec

## License

bitcask is licensed under the [MIT License](https://github.com/prologic/msgbus/blob/master/LICENSE)
