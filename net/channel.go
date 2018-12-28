package net

type Channel interface {
	WriteAndFlush(msg []byte)
	Close()
	String() string
} 
