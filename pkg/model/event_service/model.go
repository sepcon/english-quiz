package event_service

import (
	"bytes"
	"encoding/gob"
	"path"
)

type ChannelIDType = string
type MessageBodyType = []byte
type Message struct {
	Channel ChannelIDType
	Body    MessageBodyType
}

func init() {
	gob.Register(Message{})
}

func MakeChannelID(identifiers ...string) ChannelIDType {
	return path.Join("channel:", path.Join(identifiers...))
}

func (s *Message) Serialize() ([]byte, error) {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(s)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func (s *Message) Deserialize(b []byte) error {
	if s == nil {
		panic("cannot desirialize to nil object")
	}
	encoder := gob.NewDecoder(bytes.NewReader(b))
	return encoder.Decode(s)
}
