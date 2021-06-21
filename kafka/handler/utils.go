package handler

func AddToBufferedChannelIfCapacityPermits(channel chan bool, data bool) {
	if len(channel) < cap(channel) {
		channel <- data
	}
}
