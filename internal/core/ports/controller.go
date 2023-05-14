package ports

type NotifierController interface {
	MonitorDiscord() error
	MonitorAlerts() error
}
