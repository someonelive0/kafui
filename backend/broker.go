package backend

import (
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

// kafka broker
type Broker struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	ID   int    `json:"id"`
	Rack string `json:"rack"`
}

func NewBrokerFromSegmentio(broker *kafka.Broker) *Broker {
	return &Broker{
		Host: broker.Host,
		Port: broker.Port,
		ID:   broker.ID,
		Rack: broker.Rack,
	}
}

func NewBrokerArrayFromSegmentio(brokers []kafka.Broker) []Broker {
	mybrokers := make([]Broker, 0, len(brokers))
	for i := range brokers {
		mybroker := *NewBrokerFromSegmentio(&brokers[i])
		mybrokers = append(mybrokers, mybroker)
	}

	return mybrokers
}

func (p *Broker) Copy(broker *kafka.Broker) {
	p.Host = broker.Host
	p.Port = broker.Port
	p.ID = broker.ID
	p.Rack = broker.Rack
}

func (p *Broker) ToStrings() []string {
	return []string{strconv.Itoa(p.ID), p.Host, strconv.Itoa(p.Port), p.Rack}
}

func (p *Broker) AddrPort() string {
	// 	return fmt.Sprintf("%s:%d", p.Host, p.Port)
	return net.JoinHostPort(p.Host, strconv.Itoa(p.Port))
}

func BrokerHeader() []string {
	return []string{"ID", "Host", "Port", "Rack"}
}
