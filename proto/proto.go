// Package proto defines generalized protocol wireformat for the game.
package proto

type Message byte

const (
	MsgBroadcast Message = 0
	MsgOne               = 1
	MsgAll               = 2
	MsgInit              = 3
)
