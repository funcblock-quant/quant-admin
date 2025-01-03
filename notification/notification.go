package notification

type NotificationClient interface {
	SendNotification() error
}
