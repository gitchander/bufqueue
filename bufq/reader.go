package bufq

type MessageReader struct {
	p *pipe
}

func (r *MessageReader) Close() error {
	return r.CloseWithError(nil)
}

func (r *MessageReader) CloseWithError(err error) error {
	return r.p.CloseRead(err)
}

func (r *MessageReader) ReadMessage(m *Message) error {
	return r.p.ReadMessage(m)
}

func (r *MessageReader) UnreadMessage() error {
	return r.p.UnreadMessage()
}
