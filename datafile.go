package bitcask

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	pb "github.com/prologic/bitcask/proto"
	"github.com/prologic/bitcask/streampb"
)

const (
	DefaultDatafileFilename = "%09d.data"
)

var (
	ErrReadonly = errors.New("error: read only datafile")
)

type Datafile struct {
	id  int
	r   *os.File
	w   *os.File
	dec *streampb.Decoder
	enc *streampb.Encoder
}

func NewDatafile(path string, id int, readonly bool) (*Datafile, error) {
	var (
		r   *os.File
		w   *os.File
		err error
	)

	fn := filepath.Join(path, fmt.Sprintf(DefaultDatafileFilename, id))

	if !readonly {
		w, err = os.OpenFile(fn, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
		if err != nil {
			return nil, err
		}
	}

	r, err = os.Open(fn)
	if err != nil {
		return nil, err
	}

	dec := streampb.NewDecoder(r)
	enc := streampb.NewEncoder(w)

	return &Datafile{
		id:  id,
		r:   r,
		w:   w,
		dec: dec,
		enc: enc,
	}, nil
}

func (df *Datafile) Name() string {
	return df.r.Name()
}

func (df *Datafile) Close() error {
	if df.w == nil {
		return df.r.Close()
	}

	err := df.Sync()
	if err != nil {
		return err
	}
	return df.w.Close()
}

func (df *Datafile) Sync() error {
	if df.w == nil {
		return nil
	}
	return df.w.Sync()
}

func (df *Datafile) Size() (int64, error) {
	var (
		stat os.FileInfo
		err  error
	)

	if df.w == nil {
		stat, err = df.r.Stat()
	} else {
		stat, err = df.w.Stat()
	}

	if err != nil {
		return -1, err
	}

	return stat.Size(), nil
}

func (df *Datafile) Read() (pb.Entry, error) {
	var e pb.Entry
	return e, df.dec.Decode(&e)
}

func (df *Datafile) ReadAt(index int64) (e pb.Entry, err error) {
	_, err = df.r.Seek(index, os.SEEK_SET)
	if err != nil {
		return
	}
	return df.Read()
}

func (df *Datafile) Write(e pb.Entry) (int64, error) {
	if df.w == nil {
		return -1, ErrReadonly
	}

	stat, err := df.w.Stat()
	if err != nil {
		return -1, err
	}

	index := stat.Size()

	e.Index = index
	e.Timestamp = time.Now().Unix()

	err = df.enc.Encode(&e)
	if err != nil {
		return -1, err
	}

	return index, nil
}
