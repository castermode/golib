package server

import (
	"net"
	"sync"
	"time"

	log "github.com/nicholaskh/log4go"
)

const (
	CONN_TYPE_TCP = iota
	CONN_TYPE_LONG_POLLING
)

type Client struct {
	net.Conn
	LastTime    time.Time
	sessTimeout time.Duration
	Done        chan byte
	sync.Mutex
	OnClose   func()
	conn_type int8
}

func NewClient(conn net.Conn, now time.Time, sessTimeout time.Duration, ctype int8) *Client {
	return &Client{Conn: conn, LastTime: now, sessTimeout: sessTimeout, Done: make(chan byte), conn_type: ctype}
}

func (this *Client) WriteMsg(msg string) {
	this.Conn.Write([]byte(msg))
	if this.conn_type == CONN_TYPE_LONG_POLLING {
		this.Close()
	}
}

func (this *Client) CheckTimeout() {
	ticker := time.NewTicker(this.sessTimeout)
	for {
		select {
		case <-ticker.C:
			if this.IsConnected() {
				if time.Now().After(this.LastTime.Add(this.sessTimeout)) {
					log.Warn("Client connection timeout: %s", this.Conn.RemoteAddr())
					this.Close()
					return
				}
			} else {
				return
			}

		case <-this.Done:
			this.Close()
			return
		}
	}
}

func (this *Client) IsConnected() bool {
	return this.Conn != nil
}

// reentrant safe
func (this *Client) Close() {
	if this.Conn == nil {
		return
	}
	if this.OnClose != nil {
		this.OnClose()
	}
	this.Mutex.Lock()
	log.Info("Client shutdown: %s", this.Conn.RemoteAddr())
	err := this.Conn.Close()
	if err != nil {
		log.Error(err)
	}
	this.Conn = nil
	this.Mutex.Unlock()
}
