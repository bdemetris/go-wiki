package wiki

import (
	"encoding/json"
	"fmt"
	"log"
	"path"
	"runtime"
	"time"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

var open bool

func Open() error {
	var err error
	_, filename, _, _ := runtime.Caller(0) // get full path of this file
	dbfile := path.Join(path.Dir(filename), "../data.db")
	config := &bolt.Options{Timeout: 1 * time.Second}
	db, err = bolt.Open(dbfile, 0600, config)
	if err != nil {
		log.Fatal(err)
	}
	open = true
	return nil
}

func Close() {
	open = false
	db.Close()
}

func Save(p *Page) error {
	if !open {
		return fmt.Errorf("db must be opened before saving!")
	}
	err := db.Update(func(tx *bolt.Tx) error {
		pages, err := tx.CreateBucketIfNotExists([]byte("pages"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		enc, err := encode(p)
		if err != nil {
			return fmt.Errorf("could not encode Person %s: %s", p.Title, err)
		}
		err = pages.Put([]byte(p.Title), enc)
		return err
	})
	return err
}

func encode(p *Page) ([]byte, error) {
	enc, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func decode(data []byte) (*Page, error) {
	var p *Page
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func GetPage(title string) (*Page, error) {
	if !open {
		return nil, fmt.Errorf("db must be opened before saving!")
	}
	var p *Page
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte("pages"))
		k := []byte(title)
		p, err = decode(b.Get(k))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Could not get Person ID %s", title)
		return nil, err
	}
	return p, nil
}

func GetPages() ([]Page, error) {
	if !open {
		return nil, fmt.Errorf("db must be opened before saving!")
	}
	var p []Page
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte("pages"))
		b.ForEach(func(k, v []byte) error {
			var page Page
			if err := json.Unmarshal(v, &page); err != nil {
				return err
			}
			if err != nil {
				return err
			}
			p = append(p, page)
			return nil
		})
		return nil
	})
	if err != nil {
		// fmt.Printf("Could not get Person ID %s", title)
		return nil, err
	}
	return p, nil
}
