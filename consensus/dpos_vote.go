package consensus

import (
	"bytes"
	"chain/common"
	"sort"
)

const NBP = 3

var ss []Chart

type Chart struct {
	key   []byte
	value int64
}

func GetFixedBP() []string {
	return []string{
		"0x6fff74d709fd95916dcfc4d92bce760e20d23994",
		"0x80dc036328ff3c97a86af8ab3e22ab030c345bdc",
		"0x173e9cada8427ecd33ad728a91e6004fa3eacf63",
	}
}

func GetViewNumber() uint64 {
	BPs := GetFixedBP()
	return uint64(len(BPs))
}

func GetBP() []string {
	var BPs []string
	for i := 0; i < NBP; i++ {
		BPs = append(BPs, common.BytesToHexString(ss[i].key))
	}
	return BPs
}

func GetNextMaster(BPs []string, seed []byte) string {
	xsum := 1
	for _, v := range seed {
		xsum = xsum + int(v)
	}
	idx := xsum % NBP
	return BPs[idx]
}

func GetNextBFTMaster(idx uint64) string {
	BPs := GetFixedBP()
	i := idx % 2
	return BPs[i]
}

func compare(i, j int) bool {

	a := ss[i]
	b := ss[j]
	if a.value-b.value >= 0 {
		return true
	} else {
		return false
	}
}

func findVote(key []byte) int {
	idx := -1
	for i, item := range ss {
		if bytes.Equal(key, item.key) {
			idx = i
			break
		}
	}
	return idx
}

func UpdateChart(v int64, addr []byte) {
	idx := findVote(addr)
	if idx > 0 {
		ss[idx].value = v + ss[idx].value
	} else {
		ss = append(ss, Chart{
			key:   addr,
			value: v,
		})
	}
	sort.Slice(ss, compare)
}

func InitDPoS() {
}
