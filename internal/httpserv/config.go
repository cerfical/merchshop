package httpserv

import "time"

type Config struct {
	Host string
	Port string

	Timeout struct {
		ReadHeader time.Duration

		Read  time.Duration
		Write time.Duration

		Request  time.Duration
		Idle     time.Duration
		Shutdown time.Duration
	}
}
