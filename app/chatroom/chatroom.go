package chatroom

import (
	"container/list"
	"github.com/revel/revel"
	"time"
)

type EventWrapper struct {
	eventChan <-chan Event
	roomId    int
}

type Event struct {
	Type      string // "join", "leave", or "message"
	User      string
	Timestamp int    // Unix timestamp (secs)
	Text      string // What the user said (if Type == "message")
	RoomId    int
}

type SubscriptionWrapper struct {
	roomId       int
	subscription chan Subscription
}

type Subscription struct {
	RoomId int
	New    <-chan Event // New events coming in.
}

// Owner of a subscription must cancel it when they stop listening to events.
func (s Subscription) Cancel() {
	wrapper := EventWrapper{eventChan: s.New, roomId: s.RoomId}
	unsubscribe <- wrapper // Unsubscribe the channel.
	drain(s.New)           // Drain it, just in case there was a pending publish.
}

func newEvent(typ, user, msg string, roomId int) Event {
	return Event{Type: typ, User: user,
		Timestamp: int(time.Now().Unix()),
		Text:      msg, RoomId: roomId}
}

func Subscribe(roomId int) Subscription {
	revel.INFO.Println("new subscription")
	wrapper := SubscriptionWrapper{roomId: roomId, subscription: make(chan Subscription)}
	subscribe <- wrapper
	subscription := <-wrapper.subscription
	revel.INFO.Println("subscribed")

	return subscription
}

func Join(roomId int, user string) {
	revel.INFO.Println("join to room")
	publish <- newEvent("join", user, "", roomId)
}

func Say(roomId int, user, message string) {
	revel.INFO.Println("say to all")
	publish <- newEvent("message", user, message, roomId)
}

func Leave(roomId int, user string) {
	revel.INFO.Println("leave room")
	publish <- newEvent("leave", user, "", roomId)
}

var (
	subscribe   = make(chan SubscriptionWrapper, 10)
	unsubscribe = make(chan EventWrapper, 10)
	publish     = make(chan Event, 10)
)

// This function loops forever, handling the chat room pubsub
func chatroom() {
	revel.INFO.Println("chatroom start")

	subscribers := make(map[int](*list.List))

	for {
		select {
		case wrapper := <-subscribe:
			revel.INFO.Println("loop: attemp to subscribe")
			if _, ok := subscribers[wrapper.roomId]; !ok {
				subscribers[wrapper.roomId] = list.New()
			}

			subscriber := make(chan Event, 10)
			subscribers[wrapper.roomId].PushBack(subscriber)
			wrapper.subscription <- Subscription{wrapper.roomId, subscriber}

		case event := <-publish:
			revel.INFO.Println("pool: try to publish")

			for ch := subscribers[event.RoomId].Front(); ch != nil; ch = ch.Next() {
				ch.Value.(chan Event) <- event
			}

		case unsub := <-unsubscribe:
			revel.INFO.Println("unsubscribe")

			for ch := subscribers[unsub.roomId].Front(); ch != nil; ch = ch.Next() {
				if ch.Value.(chan Event) == unsub.eventChan {
					subscribers[unsub.roomId].Remove(ch)
					break
				}
			}
		}
	}
}

func init() {
	go chatroom()
}

// Drains a given channel of any messages.
func drain(ch <-chan Event) {
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		default:
			return
		}
	}
}
