package prepare

import "syscall"

const (
	lf = "singleton.lock"
)
const mask = 0644

func singleton() {
	if fd, err := syscall.Open(lf, syscall.O_CREAT|syscall.O_RDONLY, mask); err != nil {
		panic(err)
	} else if err := syscall.Flock(fd, syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		panic(err)
	}
}
