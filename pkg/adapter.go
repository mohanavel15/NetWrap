package pkg

import "github.com/songgao/water"

type Adapter struct {
	Interface *water.Interface
}

func NewAdapter(name string) (*Adapter, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = name

	ifce, err := water.New(config)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		Interface: ifce,
	}, nil
}

func (a *Adapter) TX(buffer []byte) error {
	_, err := a.Interface.Write(buffer)
	return err
}

func (a *Adapter) RX() ([]byte, error) {
	packet := make([]byte, 2000)
	n, err := a.Interface.Read(packet)
	return packet[:n], err
}
