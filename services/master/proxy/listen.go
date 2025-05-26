package proxy

func (w *Worker) listen() error {
	msg, err := w.socketListen.RecvMessage(0)
	if err != nil {
		return err
	}
	n, err := w.socketForward.SendMessage(msg)
	if err != nil {
		return err
	}
	w.Logger.Debug("Forwarded message", "size", n)
	return nil
}
