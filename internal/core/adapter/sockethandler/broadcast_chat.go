package sockethandler

import "sync"

func BroadcastMessagesHandler(mutex *sync.Mutex) func(chan Message) {
	return func(broadcast chan Message) {
		for {
			msg := <-broadcast
			for client := range Chatrooms[msg.Room] {
				if err := socketResponseAsJson(client, msg); err != nil {
					mutex.Lock()
					delete(Chatrooms[msg.Room], client)
					mutex.Unlock()
				}
			}
		}
	}
}
