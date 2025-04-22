package backend

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"

	"github.com/segmentio/kafka-go"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// kafka group of customer
type Group struct {
}

func (p *KafkaTool) ListGroups() ([]string, error) {
	client := &kafka.Client{
		Addr:      kafka.TCP(p.KafkaConfig.Brokers[0]),
		Transport: p.sharedTransport,
	}

	groupsreq := &kafka.ListGroupsRequest{Addr: client.Addr}
	groupsrep, err := client.ListGroups(context.Background(), groupsreq)
	if err != nil {
		runtime.LogErrorf(*p.Appctx, "ListGroups error: %s", err)
		return nil, err
	}
	// runtime.LogInfof(*p.Appctx, "ListGroups groups: %#v", groupsrep)

	groups := make([]string, 0, len(groupsrep.Groups))
	for i := range groupsrep.Groups {
		groups = append(groups, groupsrep.Groups[i].GroupID)
	}

	sort.Strings(groups)
	return groups, nil
}

// ####################################################################################
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
		Addr:      kafka.TCP(p.KafkaConfig.Brokers[0]),
		Transport: p.sharedTransport,
	}

	req := &kafka.DescribeGroupsRequest{
		Addr:     client.Addr,
		GroupIDs: []string{group},
	}
	resp, err := client.DescribeGroups(context.Background(), req)
	if err != nil {
		runtime.LogErrorf(*p.Appctx, "DescribeGroups error: %s", err)
		return nil, err
	}
	b, _ := json.MarshalIndent(resp, "", " ")
	runtime.LogInfof(*p.Appctx, "DescribeGroups '%s': %s\n", group, b)

	return nil, nil
}

func (p *KafkaTool) GetGroupOffset(group string) ([]GroupOffset, error) {
	client := &kafka.Client{
		Addr:      kafka.TCP(p.KafkaConfig.Brokers[0]),
		Transport: p.sharedTransport,
	}

	req := &kafka.OffsetFetchRequest{
		Addr:    client.Addr,
		GroupID: group,
	}
	resp, err := client.OffsetFetch(context.Background(), req)
	if err != nil {
		runtime.LogErrorf(*p.Appctx, "GetGroupOffset '%s' failed: %s", group, err)
		return nil, err
	}
	// b, _ := json.MarshalIndent(resp, "", " ")
	// runtime.LogInfof(*p.Appctx,"group offset %s: %s\n", group, b)

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
				runtime.LogErrorf(*p.Appctx, "GetTopicPartitionOffset '%s' failed: %s", topic, err)
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
		Addr:      kafka.TCP(p.KafkaConfig.Brokers[0]),
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
		Addr:      kafka.TCP(p.KafkaConfig.Brokers[0]),
		Transport: p.sharedTransport,
	}

	req := &kafka.DeleteGroupsRequest{
		Addr:     client.Addr,
		GroupIDs: []string{group},
	}

	resp, err := client.DeleteGroups(context.Background(), req)
	if err != nil {
		runtime.LogErrorf(*p.Appctx, "DeleteGroups '%s' failed: %s", group, err)
		return err
	}
	runtime.LogInfof(*p.Appctx, "DeleteGroups '%s': %#v", group, resp)

	if err1, ok := resp.Errors[group]; ok && err1 != nil {
		return err1
	}

	return nil
}
