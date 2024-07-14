package pkg

import (
	"fmt"
	"net"
	"runtime"
	"syscall"

	"github.com/songgao/water"
	"golang.org/x/sys/unix"
)

type Adapter struct {
	Name      string
	fd        int
	Interface *water.Interface
}

func NewAdapter(name string) (*Adapter, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}

	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		config.Name = name
	}

	ifce, err := water.New(config)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		Name:      name,
		Interface: ifce,
	}, nil
}

func (a *Adapter) SetUpFd() error {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return fmt.Errorf("Unable to get file descriptor: %s", err.Error())
	}

	ifreq, err := unix.NewIfreq(a.Name)
	if err != nil {
		return fmt.Errorf("Unable to create Ifreq: %s", err.Error())
	}

	err = unix.IoctlIfreq(fd, unix.SIOCGIFINDEX, ifreq)
	if err != nil {
		return fmt.Errorf("Unable to index interface: %s", err.Error())
	}

	a.fd = fd
	return nil
}

func (a *Adapter) SetIP(ip net.IP) error {
	ifreq, err := unix.NewIfreq(a.Name)
	if err != nil {
		return fmt.Errorf("Unable to create Ifreq: %s", err.Error())
	}

	ifreq.SetInet4Addr(ip.To4())
	err = unix.IoctlIfreq(a.fd, unix.SIOCSIFADDR, ifreq)
	if err != nil {
		return fmt.Errorf("Unable to set ip address: %s", err.Error())
	}

	return nil
}

func (a *Adapter) SetUp() error {
	ifreq, err := unix.NewIfreq(a.Name)
	if err != nil {
		return fmt.Errorf("Unable to create Ifreq: %s", err.Error())
	}

	ifreq.SetUint16(unix.IFF_UP)
	err = unix.IoctlIfreq(a.fd, unix.SIOCSIFFLAGS, ifreq)
	if err != nil {
		return fmt.Errorf("Unable to set interface to up %s", err.Error())
	}

	return nil
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

func (a *Adapter) OnError(err error) bool {
	return false
}
