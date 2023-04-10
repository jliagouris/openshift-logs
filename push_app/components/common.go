package components

// Log Parsed log
type Log struct {
	EOF   bool
	Val   map[string]interface{}
	Topic string
	Key   DataShareKey
}

type DataShareKey struct {
	ClientId string
	QueryId  uint32
	SeqNum   uint32
}
