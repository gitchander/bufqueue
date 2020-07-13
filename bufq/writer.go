package bufq

type MessageWriter struct {
	p *pipe
}

func (w *MessageWriter) Close() error {
	return w.CloseWithError(nil)
}

func (w *MessageWriter) CloseWithError(err error) error {
	return w.p.CloseWrite(err)
}

func (w *MessageWriter) WriteMessage(m *Message) error {
	return w.p.WriteMessage(m)
}
