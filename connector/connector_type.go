package connector

type NetConnectorType int

const  (
	Tcp NetConnectorType = iota + 1
	Http
)