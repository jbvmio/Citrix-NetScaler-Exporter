package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jbvmio/netscaler"
	"gopkg.in/yaml.v2"

	"github.com/dgraph-io/badger"
)

const intervalSecs = 3600

var currentMapping VIPMap

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

func (db *DB) removeLBServer(lbs lbserver) {
	db.lock.Lock()
	delete(db.lbservers, lbs.url)
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
			db.collectVIPMap2(lbs)
			log.Printf("completed vip mappings for %s\n", url)
		}
		db.setNotCollecting()
	}
}

func (db *DB) collectVIPMap2(lbs lbserver) error {
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
	currentMapping.updateMappings(lbs.url, kvMap)
	/*
		err = updateBatch(db.db, kvMap)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error updating 1 or more bindings\n")
			log.Printf("failed update for %s\n", lbs.url)
		} else {
			lbs.ready = true
			db.setLBServer(lbs)
			log.Printf("successful update for %s\n", lbs.url)
		}
	*/
	return nil
}

func (db *DB) loadVIPMap(vMap *VIPMap) {
	kvMap := make(map[string]string)
	vMap.lock.Lock()
	for _, k := range vMap.mappings {
		for a, b := range k {
			kvMap[a] = b
		}
	}
	vMap.lock.Unlock()
	err := updateBatch(db.db, kvMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error updating 1 or more bindings from file\n")
		log.Printf("error updating 1 or more bindings from file\n")
	}
}

// VIPMap contains mappings.
type VIPMap struct {
	mappings map[string]map[string]string
	lock     sync.Mutex
}

func (v *VIPMap) updateMappings(key string, maps map[string]string) {
	v.lock.Lock()
	ab, there := v.mappings[key]
	if !there {
		v.mappings[key] = make(map[string]string)
		ab = v.mappings[key]
	}
	for a, b := range maps {
		ab[a] = b
	}
	v.lock.Unlock()
}

func (v *VIPMap) exists(key string) bool {
	var there bool
	v.lock.Lock()
	_, there = v.mappings[key]
	v.lock.Unlock()
	return there
}

func (v *VIPMap) getMapping(url, key string) string {
	var val string
	v.lock.Lock()
	val = v.mappings[url][key]
	v.lock.Unlock()
	return val
}

func (v *VIPMap) getMappingYaml() (y []byte, err error) {
	v.lock.Lock()
	y, err = yaml.Marshal(v.mappings)
	v.lock.Unlock()
	return
}

func (v *VIPMap) loadMappingYaml(path string) bool {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("error loading mapping file:", err)
		return false
	}
	err = yaml.Unmarshal(b, &v.mappings)
	if err != nil {
		log.Println("error unmarshaling mapping file:", err)
		return false
	}
	if len(v.mappings) > 0 {
		return true
	}
	return false
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
