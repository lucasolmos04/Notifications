package domain

type StatusNotification string

const (
	Pending StatusNotification = "pending"
	Sent    StatusNotification = "sent"
	Failed  StatusNotification = "failed"
)
