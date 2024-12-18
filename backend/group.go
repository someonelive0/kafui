package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/segmentio/kafka-go"
)

// kafka group of customer
type Group struct {
}

// From kafka.OffsetFetchPartition
type GroupOffset struct {
	Topic           string `json:"topic"`
	Partition       int    `json:"partition"`
	FirstOffset     int64  `json:"first_offset"`
	LastOffset      int64  `json:"last_offset"`
	CommittedOffset int64  `json:"committed_offset"`

	// Consumer group metadata for this partition.
	Metadata string `json:"metadata"`
}

func (p *GroupOffset) ToStrings() []string {
	return []string{p.Topic, strconv.Itoa(p.Partition),
		strconv.FormatInt(p.FirstOffset, 10), strconv.FormatInt(p.LastOffset, 10),
		strconv.FormatInt(p.CommittedOffset, 10), p.Metadata}
}

func GroupOffsetHeader() []string {
	return []string{"Topic", "Partition", "First Offset", "Last Offset", "Committed Offset", "Metadata"}
}

// 这个函数是kafka-go库有错误，不能正常返回
func (p *KafkaTool) GetGroupDesc(group string) ([]string, error) {
	client := &kafka.Client{
		Addr:      kafka.TCP(p.kafkaConfig.Brokers[0]),
		Transport: p.sharedTransport,
	}

	req := &kafka.DescribeGroupsRequest{
		Addr:     client.Addr,
		GroupIDs: []string{group},
	}
	resp, err := client.DescribeGroups(context.Background(), req)
	if err != nil {
		return nil, err
	}
	b, _ := json.MarshalIndent(resp, "", " ")
	fmt.Printf("group descript %s: %s\n", group, b)

	return nil, nil
}

func (p *KafkaTool) GetGroupOffset(group string) ([]GroupOffset, error) {
	client := &kafka.Client{
		Addr:      kafka.TCP(p.kafkaConfig.Brokers[0]),
		Transport: p.sharedTransport,
	}

	req := &kafka.OffsetFetchRequest{
		Addr:    client.Addr,
		GroupID: group,
	}
	resp, err := client.OffsetFetch(context.Background(), req)
	if err != nil {
		return nil, err
	}
	// b, _ := json.MarshalIndent(resp, "", " ")
	// fmt.Printf("group offset %s: %s\n", group, b)

	if len(resp.Topics) == 0 {
		return nil, nil
	}

	group_offsets := make([]GroupOffset, 0)
	for topic, offsets := range resp.Topics {
		for i := range offsets {
			group_offset := GroupOffset{
				Topic:           topic,
				Partition:       offsets[i].Partition,
				CommittedOffset: offsets[i].CommittedOffset,
			}

			first, last, err := p.GetTopicPartitionOffset(topic, group_offset.Partition)
			if err != nil {
				return nil, err
			}
			group_offset.FirstOffset = first
			group_offset.LastOffset = last

			group_offsets = append(group_offsets, group_offset)
		}
	}

	return group_offsets, nil
}

// Unknown Member ID: the member id is not in the current generation, so first should GetGroupDesc success
func (p *KafkaTool) SetGroupOffset(group, topic string, partition int, offset int64) error {
	client := &kafka.Client{
		Addr:      kafka.TCP(p.kafkaConfig.Brokers[0]),
		Transport: p.sharedTransport,
	}

	req := &kafka.OffsetCommitRequest{
		Addr:     client.Addr,
		GroupID:  group,
		MemberID: group,
		Topics:   make(map[string][]kafka.OffsetCommit),
	}
	req.Topics[topic] = []kafka.OffsetCommit{{
		Partition: partition,
		Offset:    offset,
		Metadata:  "change committed offset",
	}}

	resp, err := client.OffsetCommit(context.Background(), req)
	if err != nil {
		return err
	}
	if resp.Topics != nil {
		if offsetCommitPartition, ok := resp.Topics[topic]; ok && len(offsetCommitPartition) > 0 {
			if offsetCommitPartition[0].Error != nil {
				return offsetCommitPartition[0].Error
			}
		}
	}

	return nil
}

func (p *KafkaTool) DeleteGroup(group string) error {
	client := &kafka.Client{
		Addr:      kafka.TCP(p.kafkaConfig.Brokers[0]),
		Transport: p.sharedTransport,
	}

	req := &kafka.DeleteGroupsRequest{
		Addr:     client.Addr,
		GroupIDs: []string{group},
	}

	resp, err := client.DeleteGroups(context.Background(), req)
	if err != nil {
		return err
	}
	fmt.Printf("%#v", resp)

	if err1, ok := resp.Errors[group]; ok && err1 != nil {
		return err1
	}

	return nil
}
