package kafka

import "github.com/pkg/errors"

func (s *Service) Produce(msg []byte) error {
	/*
		// add topic mapping code here

		return s.broadcastObject([]byte(msg), topic)
	*/

	// Convert topic parameter when code is added
	return s.produceObject(msg, OriginalBlockTopicFormat)
}

func (s *Service) produceObject(msg []byte, topic string) error {
	/*
		만약 추후에 메세지 종류가 많아질 경우,
		메세지 종류에 따라 topic을 매핑하여 브로드캐스트 할 수 있도록 수정

		현재는 topic을 "block"으로 고정하여 브로드캐스트
	*/
	if err := s.ProduceToTopic(s.ctx, topic, msg); err != nil {
		err := errors.Wrap(err, "could not publish message")
		return err
	}

	return nil
}
