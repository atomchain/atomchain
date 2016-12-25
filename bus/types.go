package bus

type Msg struct {
	T    string
	Data string
}

type SyncAsk struct {
	Index uint64
	Root  string
}

type SyncResp struct {
	Index     uint64
	Root      string
	Signature string
	Pubkey    string
}

type ViewSync struct {
	ViewNumber uint64
}

type ViewResp struct {
	ViewNumber uint64
}

type CommitSync struct {
	Index     uint64
	BlockData []byte
}

type CommitResp struct {
	Index uint64
	From  string
}
