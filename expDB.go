package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jbvmio/netscaler"

	"github.com/dgraph-io/badger"
)

const intervalSecs = 300

// DB handles vip mappings.
type DB struct {
	db           *badger.DB
	lbservers    map[string]lbserver
	newLBS       chan lbserver
	stopChan     chan struct{}
	isCollecting bool
	wg           sync.WaitGroup
	lock         sync.Mutex
}

func newDB(dbDir string) *DB {
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err = os.MkdirAll(dbDir, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating database dir: %v\n", err)
			os.Exit(1)
		}
	}
	db, err := badger.Open(badger.DefaultOptions(dbDir))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening database dir: %v\n", err)
		os.Exit(1)
	}

	return &DB{
		db:        db,
		lbservers: make(map[string]lbserver),
		newLBS:    make(chan lbserver, 5),
		stopChan:  make(chan struct{}),
		wg:        sync.WaitGroup{},
		lock:      sync.Mutex{},
	}
}

func (db *DB) exists(url string) (exists, ready bool) {
	db.lock.Lock()
	lbs, ok := db.lbservers[url]
	db.lock.Unlock()
	return ok, lbs.ready
}

func (db *DB) collecting() bool {
	return db.isCollecting
}

func (db *DB) setCollecting() {
	db.lock.Lock()
	db.isCollecting = true
	db.lock.Unlock()
}

func (db *DB) setNotCollecting() {
	db.lock.Lock()
	db.isCollecting = false
	db.lock.Unlock()
}

func (db *DB) setLBServer(lbs lbserver) {
	db.lock.Lock()
	db.lbservers[lbs.url] = lbs
	db.lock.Unlock()
}

func (db *DB) copy() map[string]lbserver {
	var tmp map[string]lbserver
	db.lock.Lock()
	tmp = make(map[string]lbserver, len(db.lbservers))
	for k, v := range db.lbservers {
		tmp[k] = v
	}
	db.lock.Unlock()
	return tmp
}

func (db *DB) collectAll() {
	log.Printf("starting vip mapping process ...\n")
	ticker := time.NewTicker(time.Second * intervalSecs)
collectLoop:
	for {
		select {
		case <-db.stopChan:
			log.Printf("stopping vip mapping process ...\n")
			break collectLoop
		case <-ticker.C:
			db.wg.Add(1)
			go db.collectVIPMaps(&db.wg)
		}
	}
	log.Printf("vip mapping process stopped ...\n")
}

func (db *DB) stopCollect() {
	close(db.stopChan)
	db.wg.Wait()
	db.db.Close()
}

func (db *DB) collectVIPMaps(wg *sync.WaitGroup) {
	defer wg.Done()
	switch {
	case db.collecting():
		log.Printf("vip mapping updates already in progress ...\n")
	default:
		db.setCollecting()
		mappings := db.copy()
		for url, lbs := range mappings {
			log.Printf("updating vip mappings for %s\n", url)
			db.collectVIPMap(lbs)
			log.Printf("completed vip mappings for %s\n", url)
		}
		db.setNotCollecting()
	}
}

func (db *DB) collectVIPMap(lbs lbserver) error {
	log.Printf("starting update for %s\n", lbs.url)
	nsClient, err := netscaler.NewNitroClient(lbs.url, lbs.user, lbs.pass, lbs.ignore)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating ns client: %v\n", err)
		return err
	}
	err = netscaler.Connect(nsClient)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error connecting ns client: %v\n", err)
		return err
	}
	defer func() {
		err := netscaler.Disconnect(nsClient)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error disconnecting ns client: %v\n", err)
		}
	}()
	nsBindings, err := netscaler.GetLBVSBindings(nsClient, "bulkbindings=yes")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error collecting bindings: %v\n", err)
		return err
	}
	kvMap := make(map[string]string)
	for _, b := range nsBindings.LBVServerServiceBindings {
		kvMap[b.ServiceName] = b.Name
	}
	err = updateBatch(db.db, kvMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error updating 1 or more bindings\n")
		log.Printf("failed update for %s\n", lbs.url)
	} else {
		lbs.ready = true
		db.setLBServer(lbs)
		log.Printf("successful update for %s\n", lbs.url)
	}
	return err
}

type lbserver struct {
	url    string
	user   string
	pass   string
	ignore bool
	ready  bool
}

func getValue(db *badger.DB, key string) string {
	var val []byte
	db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			//fmt.Fprintf(os.Stderr, "error getting value for key %s: %v\n", key, err)
		}
		if item == nil {
			val = []byte{}
		} else {
			val, _ = item.ValueCopy(nil)
			//fmt.Printf("%s: NOT NIL: %v\n", key, val)
		}
		return nil
	})
	return fmt.Sprintf("%s", val)
}

func updateDB(db *badger.DB, key, value []byte) error {
	return db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(key, value)
		return txn.SetEntry(e)
	})
}

func updateBatch(db *badger.DB, kv map[string]string) error {
	wb := db.NewWriteBatch()
	defer wb.Cancel()
	for k, v := range kv {
		err := wb.Set([]byte(k), []byte(v))
		if err != nil {
			log.Printf("failed to set key: %s\n", k)
		}
	}
	return wb.Flush()
}
