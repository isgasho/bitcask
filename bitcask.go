package bitcask

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultMaxDatafileSize = 1 << 20 // 1MB
)

var (
	ErrKeyNotFound = errors.New("error: key not found")
)

type Bitcask struct {
	path   string
	curr   *Datafile
	keydir *Keydir

	maxDatafileSize int64
}

func (b *Bitcask) Close() error {
	return b.curr.Close()
}

func (b *Bitcask) Sync() error {
	return b.curr.Sync()
}

func (b *Bitcask) Get(key string) ([]byte, error) {
	item, ok := b.keydir.Get(key)
	if !ok {
		return nil, ErrKeyNotFound
	}

	df, err := NewDatafile(b.path, item.FileID, true)
	if err != nil {
		return nil, err
	}
	defer df.Close()

	e, err := df.ReadAt(item.Index)
	if err != nil {
		return nil, err
	}

	return e.Value, nil
}

func (b *Bitcask) Put(key string, value []byte) error {
	index, err := b.put(key, value)
	if err != nil {
		return err
	}

	b.keydir.Add(key, b.curr.id, index, time.Now().Unix())

	return nil
}

func (b *Bitcask) Delete(key string) error {
	_, err := b.put(key, []byte{})
	if err != nil {
		return err
	}

	b.keydir.Delete(key)

	return nil
}

func (b *Bitcask) Fold(f func(key string) error) error {
	for key := range b.keydir.Keys() {
		if err := f(key); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bitcask) put(key string, value []byte) (int64, error) {
	size, err := b.curr.Size()
	if err != nil {
		return -1, err
	}

	if size >= b.maxDatafileSize {
		err := b.curr.Close()
		if err != nil {
			return -1, err
		}

		id := b.curr.id + 1
		curr, err := NewDatafile(b.path, id, false)
		if err != nil {
			return -1, err
		}
		b.curr = curr
	}

	e := NewEntry(key, value)
	return b.curr.Write(e)
}

func (b *Bitcask) setMaxDatafileSize(size int64) error {
	b.maxDatafileSize = size
	return nil
}

func MaxDatafileSize(size int64) func(*Bitcask) error {
	return func(b *Bitcask) error {
		return b.setMaxDatafileSize(size)
	}
}

func getDatafiles(path string) ([]string, error) {
	fns, err := filepath.Glob(fmt.Sprintf("%s/*.data", path))
	if err != nil {
		return nil, err
	}
	sort.Strings(fns)
	return fns, nil
}

func parseIds(fns []string) ([]int, error) {
	var ids []int
	for _, fn := range fns {
		fn = filepath.Base(fn)
		ext := filepath.Ext(fn)
		if ext != ".data" {
			continue
		}
		id, err := strconv.ParseInt(strings.TrimSuffix(fn, ext), 10, 32)
		if err != nil {
			return nil, err
		}
		ids = append(ids, int(id))
	}
	sort.Ints(ids)
	return ids, nil
}

func Merge(path string, force bool) error {
	fns, err := getDatafiles(path)
	if err != nil {
		return err
	}

	ids, err := parseIds(fns)
	if err != nil {
		return err
	}

	// Do not merge if we only have 1 Datafile
	if len(ids) <= 1 {
		return nil
	}

	// Don't merge the Active Datafile (the last one)
	fns = fns[:len(fns)-1]
	ids = ids[:len(ids)-1]

	temp, err := ioutil.TempDir("", "bitcask")
	if err != nil {
		return err
	}

	for i, fn := range fns {
		// Don't merge Datafiles whose .hint files we've already generated
		// (they are already merged); unless we set the force flag to true
		// (forcing a re-merge).
		if filepath.Ext(fn) == ".hint" && !force {
			// Already merged
			continue
		}

		id := ids[i]

		keydir := NewKeydir()

		df, err := NewDatafile(path, id, true)
		if err != nil {
			return err
		}
		defer df.Close()

		for {
			e, err := df.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}

			// Tombstone value  (deleted key)
			if len(e.Value) == 0 {
				keydir.Delete(e.Key)
				continue
			}

			keydir.Add(e.Key, ids[i], e.Index, e.Timestamp)
		}

		tempdf, err := NewDatafile(temp, id, false)
		if err != nil {
			return err
		}
		defer tempdf.Close()

		for key := range keydir.Keys() {
			item, _ := keydir.Get(key)
			e, err := df.ReadAt(item.Index)
			if err != nil {
				return err
			}

			_, err = tempdf.Write(e)
			if err != nil {
				return err
			}
		}

		err = tempdf.Close()
		if err != nil {
			return err
		}

		err = df.Close()
		if err != nil {
			return err
		}

		err = os.Rename(tempdf.Name(), df.Name())
		if err != nil {
			return err
		}

		hint := strings.TrimSuffix(df.Name(), ".data") + ".hint"
		err = keydir.Save(hint)
		if err != nil {
			return err
		}
	}

	return nil
}

func Open(path string, options ...func(*Bitcask) error) (*Bitcask, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	err := Merge(path, false)
	if err != nil {
		return nil, err
	}

	keydir := NewKeydir()

	fns, err := getDatafiles(path)
	if err != nil {
		return nil, err
	}

	ids, err := parseIds(fns)
	if err != nil {
		return nil, err
	}

	for i, fn := range fns {
		if filepath.Ext(fn) == ".hint" {
			f, err := os.Open(filepath.Join(path, fn))
			if err != nil {
				return nil, err
			}
			defer f.Close()

			hint, err := NewKeydirFromBytes(f)
			if err != nil {
				return nil, err
			}

			for key := range hint.Keys() {
				item, _ := hint.Get(key)
				keydir.Add(key, item.FileID, item.Index, item.Timestamp)
			}
		} else {
			df, err := NewDatafile(path, ids[i], true)
			if err != nil {
				return nil, err
			}

			for {
				e, err := df.Read()
				if err != nil {
					if err == io.EOF {
						break
					}
					return nil, err
				}

				// Tombstone value  (deleted key)
				if len(e.Value) == 0 {
					keydir.Delete(e.Key)
					continue
				}

				keydir.Add(e.Key, ids[i], e.Index, e.Timestamp)
			}
		}
	}

	var id int
	if len(ids) > 0 {
		id = ids[(len(ids) - 1)]
	}
	curr, err := NewDatafile(path, id, false)
	if err != nil {
		return nil, err
	}

	bitcask := &Bitcask{
		path:   path,
		curr:   curr,
		keydir: keydir,

		maxDatafileSize: DefaultMaxDatafileSize,
	}

	for _, option := range options {
		err = option(bitcask)
		if err != nil {
			return nil, err
		}
	}

	return bitcask, nil
}
