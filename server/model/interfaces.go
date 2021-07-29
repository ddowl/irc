package model

type MessageProxyStore interface {
	GetProxyMetadata() map[int]string
	AddProxy(name string) (int, error)
	GetProxy(id int) (MessageProxy, error)
	DeleteProxy(id int) error
}

type MessageProxy interface {
	Subscribable
	Broadcaster
}

type Subscribable interface {
	Join(tag string, callbackUrl string) error
	Leave(tag string) error
}

type Broadcaster interface {
	PostMessage(tag string, message string) error
}
