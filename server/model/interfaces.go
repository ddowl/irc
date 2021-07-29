package model

type MessageProxyStore interface {
	GetMetadata() []ProxyMetadata
	AddProxy(name string) (int, error)
	GetProxy(id int) (MessageProxy, error)
	DeleteProxy(id int) error
}

type MessageProxy interface {
	GetMetadata() *ProxyMetadata
	Subscribable
	Broadcaster
}

type Subscribable interface {
	Join(tag string, callbackUrl string) error
	Leave(tag string) error
	HasJoined(tag string) bool
}

type Broadcaster interface {
	PostMessage(tag string, message string) error
}

type ProxyMetadata struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
