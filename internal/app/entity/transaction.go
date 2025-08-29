package entity

type IsolationLevel uint8

const (
	RepeatableRead IsolationLevel = iota + 1
	Serializable
)
