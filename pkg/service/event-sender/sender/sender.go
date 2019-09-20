package sender

import "log"

func Send(message []byte) {
	log.Printf("Hey hey! Don't foget about %v", string(message))
}
