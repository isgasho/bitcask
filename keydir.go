package bitcask

import (
	"bytes"
	"encoding/gob"
	"io"
	"io/ioutil"
	"sync"
)

type Item struct {
	FileID    int
	Index     int64
	Timestamp int64
}

type Keydir struct {
	sync.RWMutex
	kv map[string]Item
}

func NewKeydir() *Keydir {
	return &Keydir{
		kv: make(map[string]Item),
	}
}

func (k *Keydir) Add(key string, fileid int, index, timestamp int64) {
	k.Lock()
	defer k.Unlock()

	k.kv[key] = Item{
		FileID:    fileid,
		Index:     index,
		Timestamp: timestamp,
	}
}

func (k *Keydir) Get(key string) (Item, bool) {
	k.RLock()
	defer k.RUnlock()

	item, ok := k.kv[key]
	return item, ok
}

func (k *Keydir) Delete(key string) {
	k.Lock()
	defer k.Unlock()

	delete(k.kv, key)
}

func (k *Keydir) Keys() chan string {
	ch := make(chan string)
	go func() {
		for k := range k.kv {
			ch <- k
		}
		close(ch)
	}()
	return ch
}

func (k *Keydir) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(k.kv)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (k *Keydir) Save(fn string) error {
	data, err := k.Bytes()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fn, data, 0644)
}

func NewKeydirFromBytes(r io.Reader) (*Keydir, error) {
	k := NewKeydir()
	dec := gob.NewDecoder(r)
	err := dec.Decode(&k.kv)
	if err != nil {
		return nil, err
	}
	return k, nil
}
