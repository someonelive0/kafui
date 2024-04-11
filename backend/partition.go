package backend

import "strconv"

// Copy from segmentio/kafka-go
type Partition struct {
	// Name of the topic that the partition belongs to, and its index in the
	// topic.
	Topic       string `json:"topic"`
	ID          int    `json:"id"`
	FirstOffset int64  `json:"first_offset"`
	LastOffset  int64  `json:"last_offset"`
	Number      int64  `json:"number"` // number = LastOffset - FirstOffset

	// Leader, replicas, and ISR for the partition.
	//
	// When no physical host is known to be running a broker, the Host and Port
	// fields will be set to the zero values. The logical broker ID is always
	// set to the value known to the kafka cluster, even if the broker is not
	// currently backed by a physical host.
	Leader   Broker   `json:"leader"`
	Replicas []Broker `json:"replicas"`
	Isr      []Broker `json:"isr"`

	// Available only with metadata API level >= 6:
	OfflineReplicas []Broker `json:"offline_replicas"`
}

func (p *Partition) ToStrings() []string {
	return []string{p.Topic, strconv.Itoa(p.ID), strconv.Itoa(int(p.FirstOffset)), strconv.Itoa(int(p.LastOffset)), p.Leader.AddrPort()}
}

func PartitionrHeader() []string {
	return []string{"Topic", "ID", "First", "Last", "Leader"}
}
