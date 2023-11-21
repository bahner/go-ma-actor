package space

import log "github.com/sirupsen/logrus"

// A way to easiy close a space.
func (s *Space) Close() {

	if s.Public != nil {
		err := s.Public.Close()
		if err != nil {
			log.Warnf("Failed to close keyAgreement topic: %v", err)
		}
	}

	if s.Private != nil {
		err := s.Private.Close()
		if err != nil {
			log.Warnf("Failed to close keyAgreement topic: %v", err)
		}
	}

}
