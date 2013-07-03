// tcpchan project tcpchan.go
package tcpchan

import (
  "encoding/gob"
	"net"
)

// Wrapper for the payload (Value) so gob can serialize it
type data struct {
	Value interface{}
}

// Creates a new channel listening on addr
func NewListener(addr string) (chan interface{}, error) {
	serv, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	ch := make(chan interface{})
	go listen(serv, ch)

	return ch, nil
}

// Creates a new channel connecting to addr
func NewDialer(addr string) (chan interface{}, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	ch := make(chan interface{})
	go write(conn, ch)

	return ch, nil
}

// Listenes on a socket and starts writer for the first connection
func listen(serv net.Listener, ch chan interface{}) {
	conn, err := serv.Accept()
	if err != nil {
		close(ch)
		return
	}
	serv.Close()

	write(conn, ch)
}

// Reads data from the connection conn and writes it into ch
func read(conn net.Conn, ch chan interface{}) {
	defer close(ch) // make sure we close the channel when we stop reading

	dec := gob.NewDecoder(conn)
	buf := data{}
	for {
		err := dec.Decode(&buf) // try to de-serialize the data
		if err != nil {
			return
		}
		ch <- buf.Value // unpack the payload
	}
}

// Handles writing to the remote channel and incoming data
func write(conn net.Conn, ch chan interface{}) {
	defer conn.Close() // make sure we close the connection when we are done writing

	in := make(chan interface{}, 256) // We need to buffer incoming data
	go read(conn, in)                 // start the reader so we can use select

	enc := gob.NewEncoder(conn)

	for {
		select {
		case i, ok := <-ch: // seems we want to send data
			if !ok {
				return
			}
			enc.Encode(data{i})
		case i, ok := <-in: // seems we received data
			if !ok {
				close(ch)
				return
			}
			for {
				select {
				case ch <- i:
					break
				case i, ok := <-ch: // seems we want to send data (so we are unable to write the received data)
					if !ok {
						return
					}
					enc.Encode(data{i})
				}
			}
		}
	}
}
