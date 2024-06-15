package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/tidwall/resp"
)

type Peer struct {
	conn  net.Conn
	msgCh chan Message
	delCh chan *Peer
}

func (p *Peer) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}

func NewPeer(conn net.Conn, msgCh chan Message, delCh chan *Peer) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
		delCh: delCh,
	}
}

func (p *Peer) readLoop() error {
	rd := resp.NewReader(p.conn)

	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			p.delCh <- p
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Printf("Read %s\n", v.Type())

		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				switch value.String() {
				case CommandGET:
					if len(v.Array()) != 2 {
						return fmt.Errorf("invalid number of arguments for GET command")
					}
					cmd := GetCommand{
						key: v.Array()[1].Bytes(),
					}

					p.msgCh <- Message{
						cmd:  cmd,
						peer: p,
					}

				case CommandSET:
					if len(v.Array()) != 3 {
						return fmt.Errorf("invalid number of arguments for SET command")
					}
					cmd := SetCommand{
						key: v.Array()[1].Bytes(),
						val: v.Array()[2].Bytes(),
					}

					p.msgCh <- Message{
						cmd:  cmd,
						peer: p,
					}
				}
			}
		}
	}
	return nil
}