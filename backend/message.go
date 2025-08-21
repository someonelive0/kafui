package backend

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Message struct {
	Topic         string   `json:"topic"`
	Partition     int      `json:"partition"` // Partition is read-only and MUST NOT be set when writing messages
	Offset        int64    `json:"offset"`
	HighWaterMark int64    `json:"high_water_mark"`
	Key           string   `json:"key"`
	Value         string   `json:"value"`
	Headers       []Header `json:"headers"`
	Time          string   `json:"time"` // time of writing the message
}

type Header struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
}

func NewMessageFromSegmentio(m *kafka.Message) *Message {

	msg := &Message{
		Topic:         m.Topic,
		Partition:     m.Partition,
		Offset:        m.Offset,
		HighWaterMark: m.HighWaterMark,
		Key:           string(m.Key),
		Value:         string(m.Value),
		Headers:       make([]Header, 0, len(m.Headers)),
		// Time:          m.Time.In(loc).Format(time.RFC3339),
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil { // win10 will cause this error: The device is not ready
		// fmt.Printf("time.LoadLocation failed %s\n", err)
		msg.Time = m.Time.Format(time.RFC3339)
	} else {
		msg.Time = m.Time.In(loc).Format(time.RFC3339)
	}

	for i := range m.Headers {
		header := Header{
			Key:   m.Headers[i].Key,
			Value: m.Headers[i].Value,
		}
		msg.Headers = append(msg.Headers, header)
	}

	return msg
}

// partition = -1 means not set partition
func (p *KafkaTool) ReadMsgs2Ch(ctx context.Context, topic string, partition, limit int, ch chan *Message) error {
	// fire read partition first and last offset
	firstOffet, lastOffset, err := p.GetTopicPartitionOffset(topic, partition)
	if err != nil {
		runtime.LogErrorf(*p.Appctx, "GetTopicPartitionOffset error: %s", err)
		return err
	}
	if lastOffset-firstOffet == 0 { // msg bumber is 0 menas topic is empty
		return nil
	}

	rconfig := kafka.ReaderConfig{
		Brokers: p.KafkaConfig.Brokers,
		// GroupID:  KafkaConfig.Group, // 指定消费者组id
		Topic:            topic,
		Dialer:           p.dialer,
		MaxBytes:         10e6, // 10MB
		ReadBatchTimeout: time.Second * 5,
		SessionTimeout:   time.Second * 5,
	}
	if partition > -1 {
		rconfig.Partition = partition
	}
	r := kafka.NewReader(rconfig)
	// defer r.Close() // will cause 6 second

	// set beginOffet when limit < 0
	if limit < 0 {
		if lastOffset+int64(limit) > firstOffet {
			if err = r.SetOffset(lastOffset + int64(limit)); err != nil {
				runtime.LogErrorf(*p.Appctx, "SetOffset on topic '%s %v' error: %s", topic, partition, err)
				return nil
			}
		}
		limit = -1 * limit
	}

	// ctx := context.Background()
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	for i := 0; limit == 0 || i < limit; i++ {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				// fmt.Printf("kafka get context canceled:%v\n", err)
				break
			}
			if errors.Is(err, io.EOF) { // 当reader.Close后，进入这个分支, but reach lastOffset not work
				// fmt.Printf("kafka get eof")
				break
			}

			runtime.LogErrorf(*p.Appctx, "FetchMessage error: %s", err)
			break
		}

		// fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %d\n", m.Topic, m.Partition, m.Offset, string(m.Key), len(m.Value))
		msg := NewMessageFromSegmentio(&m)
		ch <- msg

		// When read EOF return, that m.Offset == lastOffset-1
		if m.Offset == lastOffset-1 {
			// fmt.Printf("kafka get EOF and return: m.Offet = last-1: %v=%v", m.Offset, lastOffset)
			break
		}
	}

	go r.Close()
	return nil
}

// func (p *KafkaTool) ReadMsgs1(topic string, partition int, timeout int) ([]Message, error) {
// 	rconfig := kafka.ReaderConfig{
// 		Brokers: p.KafkaConfig.Brokers,
// 		// GroupID:  KafkaConfig.Group, // 指定消费者组id
// 		Topic:            topic,
// 		Dialer:           p.dialer,
// 		MaxBytes:         10e6, // 10MB
// 		ReadBatchTimeout: time.Second * 5,
// 		SessionTimeout:   time.Second * 5,
// 	}
// 	if partition > -1 {
// 		rconfig.Partition = partition
// 	}
// 	r := kafka.NewReader(rconfig)
// 	defer r.Close() // will cause 6 second

// 	// ctx := context.Background()
// 	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
// 	count := 0
// 	msgs := make([]Message, 0)
// 	for {
// 		m, err := r.FetchMessage(ctx)
// 		if err != nil {
// 			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
// 				// fmt.Printf("kafka get context canceled:%v", err)
// 				break
// 			}
// 			if errors.Is(err, io.EOF) { // 当reader.Close后，进入这个分支
// 				// fmt.Printf("kafka get eof")
// 				break
// 			}
// 			// runtime.LogErrorf(*p.Appctx, "FetchMessage error: %s", err)
// 			break
// 		}

// 		// fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %d", m.Topic, m.Partition, m.Offset, string(m.Key), len(m.Value))
// 		msg := NewMessageFromSegmentio(&m)
// 		count++
// 		msgs = append(msgs, *msg)
// 	}

// 	return msgs, nil
// }

func (p *KafkaTool) ReadMsgs(topic string, partition int, timeout int) ([]Message, error) {
	return p.ReadMsgsLimit(topic, partition, -1, timeout)
}

func (p *KafkaTool) ReadMsgsLimit(topic string, partition, limit, timeout int) ([]Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	ch := make(chan *Message, 10)

	go func() {
		defer close(ch)
		defer cancel()
		err := p.ReadMsgs2Ch(ctx, topic, partition, limit, ch)
		if err != nil {
			runtime.LogErrorf(*p.Appctx, "GetMessages failed: %s", err)
			return
		}
	}()

	count := 0
	msgs := make([]Message, 0)
	for msg := range ch {
		// fmt.Printf("%d > of %s %d : %s = %s\n", count,
		// 	msg.Time.Format(time.DateTime), msg.Partition, string(msg.Key), string(msg.Value))
		count++
		msgs = append(msgs, *msg)
	}
	// fmt.Printf("total get %d messages of topic %s\n", count, topic)

	return msgs, nil
}

func (p *KafkaTool) WriteMsg(topic string, key, value string) error {
	w := &kafka.Writer{
		Addr:         kafka.TCP(p.KafkaConfig.Brokers...),
		Topic:        topic,
		Transport:    p.sharedTransport,
		Balancer:     &kafka.LeastBytes{}, // 指定分区的balancer模式为最小字节分布
		RequiredAcks: kafka.RequireAll,    // ack模式
		Async:        false,               // 同步
		WriteTimeout: time.Second * 5,
	}
	defer w.Close()

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(value)},
	)
	if err != nil {
		runtime.LogErrorf(*p.Appctx, "WriteMsg error: %s", err)
		return err
	}

	if err = w.Close(); err != nil {
		return err
	}

	return nil
}
