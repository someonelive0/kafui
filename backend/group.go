package backend

import (
	"context"
	"encoding/json"
	"fmt"
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
/*
kafka-go v0.4.48 Now work fine, return
{
 "Error": null,
 "GroupID": "testgroup",
 "GroupState": "Empty",
 "Members": null
}
*/
func (p *KafkaTool) GetGroupDesc(group string) ([]byte, error) {
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
		if p.Appctx != nil {
			runtime.LogErrorf(*p.Appctx, "DescribeGroups error: %s", err)
		}
		return nil, err
	}
	if len(resp.Groups) == 0 {
		return nil, fmt.Errorf("not found group '%s'", group)
	}
	b, _ := json.MarshalIndent(resp.Groups[0], "", " ")
	if p.Appctx != nil {
		runtime.LogInfof(*p.Appctx, "DescribeGroups '%s': %s\n", group, b)
	}

	return b, nil
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
func (p *KafkaTool) SetGroupOffset(group, topic string, partition int, new_offset int64) error {
	group_offsets, err := p.GetGroupOffset(group)
	if err != nil {
		return err
	}

	var my_group_offset *GroupOffset
	for _, group_offset := range group_offsets {
		if group_offset.Topic == topic && group_offset.Partition == partition {
			my_group_offset = &group_offset
			break
		}
	}
	if my_group_offset == nil {
		return fmt.Errorf("not found group offset: %s:%d", topic, partition)
	}
	if my_group_offset.CommittedOffset == new_offset { // no need to commit
		return nil
	}
	if new_offset < my_group_offset.FirstOffset {
		return fmt.Errorf("new offset %d is less than first offset %d", new_offset, my_group_offset.FirstOffset)
	}
	if new_offset > my_group_offset.LastOffset {
		return fmt.Errorf("new offset %d is bigger than last offset %d", new_offset, my_group_offset.LastOffset)
	}

	consumergroup, err := kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
		ID:      group,
		Brokers: p.KafkaConfig.Brokers,
		Topics:  []string{topic},
	})
	if err != nil {
		// fmt.Printf("error creating consumer group: %+v\n", err)
		return err
	}
	defer consumergroup.Close()

	gen, err := consumergroup.Next(context.TODO())
	if err != nil {
		// fmt.Printf("error getting next generation: %+v\n", err)
		return err
	}

	// assignments is empty.
	// assignments := gen.Assignments[topic]
	// for _, assignment := range assignments {
	// 	assign_partition, assign_offset := assignment.ID, assignment.Offset
	// 	fmt.Printf("topic %s: partition: %d, offset: %d\n", topic, assign_partition, assign_offset)
	// }

	err = gen.CommitOffsets(map[string]map[int]int64{
		topic: {
			partition: new_offset, // the offset to commit
		},
	})
	if err != nil {
		// fmt.Printf("error committing offsets next generation: %+v\n", err)
		return err
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
