package bitcask

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	var (
		db      *Bitcask
		testdir string
		err     error
	)

	assert := assert.New(t)

	testdir, err = ioutil.TempDir("", "bitcask")
	assert.NoError(err)

	t.Run("Open", func(t *testing.T) {
		db, err = Open(testdir)
		assert.NoError(err)
	})

	t.Run("Put", func(t *testing.T) {
		err = db.Put("foo", []byte("bar"))
		assert.NoError(err)
	})

	t.Run("Get", func(t *testing.T) {
		val, err := db.Get("foo")
		assert.NoError(err)
		assert.Equal([]byte("bar"), val)
	})

	t.Run("Delete", func(t *testing.T) {
		err := db.Delete("foo")
		assert.NoError(err)
		_, err = db.Get("foo")
		assert.Error(err)
		assert.Equal(err.Error(), "error: key not found")
	})

	t.Run("Sync", func(t *testing.T) {
		err = db.Sync()
		assert.NoError(err)
	})

	t.Run("Close", func(t *testing.T) {
		err = db.Close()
		assert.NoError(err)
	})
}

func TestDeletedKeys(t *testing.T) {
	assert := assert.New(t)

	testdir, err := ioutil.TempDir("", "bitcask")
	assert.NoError(err)

	t.Run("Setup", func(t *testing.T) {
		var (
			db  *Bitcask
			err error
		)

		t.Run("Open", func(t *testing.T) {
			db, err = Open(testdir)
			assert.NoError(err)
		})

		t.Run("Put", func(t *testing.T) {
			err = db.Put("foo", []byte("bar"))
			assert.NoError(err)
		})

		t.Run("Get", func(t *testing.T) {
			val, err := db.Get("foo")
			assert.NoError(err)
			assert.Equal([]byte("bar"), val)
		})

		t.Run("Delete", func(t *testing.T) {
			err := db.Delete("foo")
			assert.NoError(err)
			_, err = db.Get("foo")
			assert.Error(err)
			assert.Equal("error: key not found", err.Error())
		})

		t.Run("Sync", func(t *testing.T) {
			err = db.Sync()
			assert.NoError(err)
		})

		t.Run("Close", func(t *testing.T) {
			err = db.Close()
			assert.NoError(err)
		})
	})

	t.Run("Reopen", func(t *testing.T) {
		var (
			db  *Bitcask
			err error
		)

		t.Run("Open", func(t *testing.T) {
			db, err = Open(testdir)
			assert.NoError(err)
		})

		t.Run("Get", func(t *testing.T) {
			_, err = db.Get("foo")
			assert.Error(err)
			assert.Equal("error: key not found", err.Error())
		})

		t.Run("Close", func(t *testing.T) {
			err = db.Close()
			assert.NoError(err)
		})
	})
}

func TestMerge(t *testing.T) {
	assert := assert.New(t)

	testdir, err := ioutil.TempDir("", "bitcask")
	assert.NoError(err)

	t.Run("Setup", func(t *testing.T) {
		var (
			db  *Bitcask
			err error
		)

		t.Run("Open", func(t *testing.T) {
			db, err = Open(testdir, MaxDatafileSize(1024))
			assert.NoError(err)
		})

		t.Run("Put", func(t *testing.T) {
			for i := 0; i < 1024; i++ {
				err = db.Put(string(i), []byte(strings.Repeat(" ", 1024)))
				assert.NoError(err)
			}
		})

		t.Run("Get", func(t *testing.T) {
			for i := 0; i < 32; i++ {
				err = db.Put(string(i), []byte(strings.Repeat(" ", 1024)))
				assert.NoError(err)
				val, err := db.Get(string(i))
				assert.NoError(err)
				assert.Equal([]byte(strings.Repeat(" ", 1024)), val)
			}
		})

		t.Run("Sync", func(t *testing.T) {
			err = db.Sync()
			assert.NoError(err)
		})

		t.Run("Close", func(t *testing.T) {
			err = db.Close()
			assert.NoError(err)
		})
	})

	t.Run("Merge", func(t *testing.T) {
		var (
			db  *Bitcask
			err error
		)

		t.Run("Open", func(t *testing.T) {
			db, err = Open(testdir)
			assert.NoError(err)
		})

		t.Run("Get", func(t *testing.T) {
			for i := 0; i < 32; i++ {
				val, err := db.Get(string(i))
				assert.NoError(err)
				assert.Equal([]byte(strings.Repeat(" ", 1024)), val)
			}
		})

		t.Run("Close", func(t *testing.T) {
			err = db.Close()
			assert.NoError(err)
		})
	})
}

func BenchmarkGet(b *testing.B) {
	testdir, err := ioutil.TempDir("", "bitcask")
	if err != nil {
		b.Fatal(err)
	}

	db, err := Open(testdir)
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	err = db.Put("foo", []byte("bar"))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val, err := db.Get("foo")
		if err != nil {
			b.Fatal(err)
		}
		if string(val) != "bar" {
			b.Errorf("expected val=bar got=%s", val)
		}
	}
}

func BenchmarkPut(b *testing.B) {
	testdir, err := ioutil.TempDir("", "bitcask")
	if err != nil {
		b.Fatal(err)
	}

	db, err := Open(testdir)
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := db.Put(fmt.Sprintf("key%d", i), []byte("bar"))
		if err != nil {
			b.Fatal(err)
		}
	}
}
