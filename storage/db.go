package storage

import (
	"bytes"
	"chain/consensus/block"
	"chain/storage/grockdb"
	"encoding/json"
	"sync"

	"github.com/cbergoon/merkletree"
	"github.com/golang/groupcache/lru"
)

type DB struct {
	rdb    *gorocksdb.DB
	wo     *gorocksdb.WriteOptions
	ro     *gorocksdb.ReadOptions
	mu     sync.Mutex
	caches *lru.Cache
}

func NewRocksDB(dir string) *gorocksdb.DB {
	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	rdb, err := gorocksdb.OpenDb(opts, dir)
	if err != nil {
		panic(err)
	}

	return rdb
}

func New(rdb *gorocksdb.DB) *DB {
	db := &DB{rdb: rdb}
	db.wo = gorocksdb.NewDefaultWriteOptions()
	db.ro = gorocksdb.NewDefaultReadOptions()
	db.caches = lru.New(1000)
	db.RawSet([]byte{MAXBYTE}, nil)
	return db
}

func (d *DB) objFromCache(key []byte, e ElementType) interface{} {
	d.mu.Lock()
	defer d.mu.Unlock()

	skey := string(key)
	obj, ok := d.caches.Get(skey)
	if !ok {
		switch e {
		case HASH:
			obj = NewHashElement(d, key)
		case LIST:
			obj = NewListElement(d, key)
		case SORTEDSET:
			obj = NewSortedSetElement(d, key)
		}
		d.caches.Add(skey, obj)
	}
	return obj
}

func (d *DB) List(key []byte) *ListElement {
	return d.objFromCache(key, LIST).(*ListElement)
}

func (d *DB) FLushAll() {
	// delete all
}

func (d *DB) Keys() []string {
	keyList := []string{}
	batch := gorocksdb.NewWriteBatch()
	defer batch.Destroy()

	d.PrefixEnumerate(KEY, IterForward, func(i int, key, value []byte, quit *bool) {
		keyName, _ := SplitKeyName(key)
		keyList = append(keyList, keyName)
	})

	return keyList
}

func (d *DB) Delete(key []byte) error {
	return d.RawDelete(key)
}

func (d *DB) TypeOf(key []byte) ElementType {
	c := ElementType(NONE)
	prefix := bytes.Join([][]byte{KEY, key, SEP}, nil)
	d.PrefixEnumerate(prefix, IterForward, func(i int, key, value []byte, quit *bool) {
		c = ElementType(key[len(prefix):][0])
		*quit = true
	})
	return c
}

func (d *DB) Get(key []byte) ([]byte, error) {
	return d.RawGet(rawKey(key, STRING))
}

func (d *DB) GetList(key []byte) ([]byte, error) {
	return d.RawGet(rawKey(key, LIST))
}

func (d *DB) Set(key, value []byte) error {
	return d.RawSet(rawKey(key, STRING), value)
}

func (d *DB) WriteBatch(batch *gorocksdb.WriteBatch) error {
	return d.rdb.Write(d.wo, batch)
}

func (d *DB) RawGet(key []byte) ([]byte, error) {
	return d.rdb.GetBytes(d.ro, key)
}

func (d *DB) RawSet(key, value []byte) error {
	return d.rdb.Put(d.wo, key, value)
}

func (d *DB) RawDelete(key []byte) error {
	return d.rdb.Delete(d.wo, key)
}

func (d *DB) Close() {
	d.wo.Destroy()
	d.ro.Destroy()
	d.rdb.Close()
}

func (d *DB) PrefixEnumerate(prefix []byte, direction IterDirection, fn func(i int, key, value []byte, quit *bool)) {
	min := prefix
	max := append(prefix, MAXBYTE)
	j := -1
	d.RangeEnumerate(min, max, direction, func(i int, key, value []byte, quit *bool) {
		if bytes.HasPrefix(key, prefix) {
			j++
			fn(j, key, value, quit)
		} else {
			*quit = true
		}
	})
	return
}

func (d *DB) RangeEnumerate(min, max []byte, direction IterDirection, fn func(i int, key, value []byte, quit *bool)) {
	opts := gorocksdb.NewDefaultReadOptions()
	opts.SetFillCache(false)
	defer opts.Destroy()
	iter := d.rdb.NewIterator(opts)
	defer iter.Close()
	d.Enumerate(iter, min, max, direction, fn)
}

func (d *DB) GetBlockIndex(ns string) uint64 {
	dl := d.List([]byte(ns))
	idx := dl.Len()
	return uint64(idx)
}

func (d *DB) SaveBlock(block *block.Block, ns string) {
	data, _ := block.Marshal()
	dl := d.List([]byte(ns))
	dl.RPush(data)
}

func UnmarshalBlock(data []byte) *block.Block {
	b := block.Block{}
	bf := bytes.NewBuffer(data)
	dec := json.NewDecoder(bf)
	if err := dec.Decode(&b); err != nil {
		return nil
	}
	return &b
}

type LeafContent struct {
	x []byte
}

func (l LeafContent) CalculateHash() ([]byte, error) {
	return l.x, nil
}

func (l LeafContent) Equals(other merkletree.Content) (bool, error) {
	return bytes.Equal(l.x, other.(LeafContent).x), nil
}

func (d *DB) GetBlockMT(ns string) []byte {
	blockList := d.List([]byte(ns))
	lastIndex := uint64(blockList.Len())
	if lastIndex == 0 {
		return nil
	}
	var hashList [][]byte
	for i := uint64(0); i < lastIndex; i++ {
		blkData, _ := blockList.Index(int64(i))
		blk := UnmarshalBlock(blkData)
		hash, _ := blk.Hash()
		hashList = append(hashList, hash)
	}
	l := len(hashList)
	var list []merkletree.Content
	for index := 0; index < l; index++ {
		aleaf := LeafContent{
			x: hashList[index],
		}
		list = append(list, aleaf)
	}
	t, _ := merkletree.NewTree(list)
	return t.MerkleRoot()
}

func (d *DB) Enumerate(iter *gorocksdb.Iterator, min, max []byte, direction IterDirection, fn func(i int, key, value []byte, quit *bool)) {
	found := false
	if direction == IterBackward {
		if len(max) == 0 {
			iter.SeekToLast()
		} else {
			iter.Seek(max)
		}
	} else {
		if len(min) == 0 {
			iter.SeekToFirst()
		} else {
			iter.Seek(min)
		}

	}
	found = iter.Valid()
	if !found {
		return
	}

	i := -1
	if found && between(iter.Key().Data(), min, max) {
		i++
		quit := false
		fn(i, iter.Key().Data(), iter.Value().Data(), &quit)
		if quit {
			return
		}
	}
	for {
		found = false
		if direction == IterBackward {
			iter.Prev()
		} else {
			iter.Next()
		}
		found = iter.Valid()
		if found && between(iter.Key().Data(), min, max) {
			i++
			quit := false
			fn(i, iter.Key().Data(), iter.Value().Data(), &quit)
			if quit {
				return
			}
		} else {
			break
		}
	}

	return
}
