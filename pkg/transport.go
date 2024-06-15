package pkg

import "log"

type Transport interface {
	TX(buffer []byte) error
	RX() ([]byte, error)
}

func Relay(src, dst Transport) {
	for {
		buffer, err := src.RX()
		if err != nil {
			log.Println("Unable Read: ", err.Error())
			continue
		}

		err = dst.TX(buffer)
		if err != nil {
			log.Println("Unable Write: ", err.Error())
			continue
		}
	}
}
