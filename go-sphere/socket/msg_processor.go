package socket

func (s *Service) ProcessMsg(msgData []byte) error {
	if err := s.cfg.Kafka.Produce(msgData); err != nil {
		log.Error("Kafka Produce Error:", err)
		return err
	}

	return nil
}
