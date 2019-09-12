package main

import (
	"fmt"
	"strings"
	"time"
)

type chatter interface {
	sendMsg(receiver entity, data string)
	chat()
	generateMsg() string
}

type message struct {
	src  *entity
	dest *entity
	data string
}

type entity struct {
	inputChannel    <-chan *message
	outputChannel   chan<- *message
	inputBandwidth  time.Duration
	outputBandwidth time.Duration
	name            string
	contacts        []*entity
}

func NewEntity(name string, inputBandwidth, outputBandwidth time.Duration) *entity {
	e := new(entity)
	e.inputChannel = make(<-chan *message)
	e.outputChannel = make(chan<- *message)
	e.inputBandwidth = inputBandwidth
	e.outputBandwidth = outputBandwidth
	e.name = name
	return e
}

func (e *entity) addContact(anotherEntity *entity) {
	e.contacts = append(e.contacts, anotherEntity)
}

func exchangeContacts(e1, e2 *entity) {
	e1.addContact(e2)
	e2.addContact(e1)
}

func (e *entity) chat() {
	fmt.Println("Starting Channel: " + e.name)
	for {
		select {
		case receivedMsg := <-e.inputChannel:
			fmt.Println(strings.Join([]string{e.name, "received:", receivedMsg.data, receivedMsg.src.name}, ""))
			e.sendMsg(e.contacts[0], "Received: "+receivedMsg.data)
			time.Sleep(e.inputBandwidth)
		default:
			time.Sleep(2 * e.outputBandwidth)
			fmt.Println(e.name + " sent msg to " + e.contacts[0].name)
			e.sendMsg(e.contacts[0], e.generateMsg())
		}
	}

}

func (e *entity) generateMsg() string {
	return "random string"
}

func (e *entity) sendMsg(receiver *entity, data string) {
	select {
	case e.outputChannel <- &message{e, receiver, data}:
		fmt.Println(strings.Join([]string{e.name, "sent:", data, receiver.name}, ""))
		time.Sleep(e.outputBandwidth)
	case <-time.After(2 * e.outputBandwidth):
		fmt.Println(e.name + " send msg time out")
	}
}

func main() {
	patrick := NewEntity("Patrick", time.Second, time.Second)
	ruth := NewEntity("Ruth", time.Second, time.Second)
	exchangeContacts(patrick, ruth)
	go patrick.chat()
	go ruth.chat()
	time.Sleep(10 * time.Second)
}
