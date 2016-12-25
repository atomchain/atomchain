package net

type Pool struct {
	txs [][]byte
}

func NewPool() *Pool {
	pool := Pool{}
	return &pool
}

func PushTx(pool *Pool, tx []byte) {
	pool.txs = append(pool.txs, tx)
}

func PoolSize(pool *Pool) int {
	return len(pool.txs)
}

func GetTxs(pool *Pool) [][]byte {
	return pool.txs
}

func EmptyTxs(pool *Pool) {
	pool.txs = pool.txs[:0]
}
