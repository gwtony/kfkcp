package kfkc
import (
	"os"
	"errors"
	"syscall"
	"os/signal"
)


func Daemon(nochdir int) error{
	signal.Ignore(syscall.SIGINT, syscall.SIGHUP, syscall.SIGPIPE)

	if syscall.Getppid() == 1 {
		return nil
	}

	ret1, ret2, eno := syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
	if eno != 0 || ret2 < 0 {
		return errors.New("Fork fail")
	}

	/* exit parent */
	if ret1 > 0 {
		os.Exit(0)
	}

	_ = syscall.Umask(0)

	ret, err := syscall.Setsid()
	if ret < 0 || err != nil {
		return errors.New("Set sid failed")
	}

	if nochdir == 0 {
		err = os.Chdir("/")
		if err != nil {
			return errors.New("Chdir failed")
		}
	}

	file, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
	if err == nil {
		fd := file.Fd()
		syscall.Dup2(int(fd), int(os.Stdin.Fd()))
		syscall.Dup2(int(fd), int(os.Stdout.Fd()))
	//	syscall.Dup2(int(fd), int(os.Stderr.Fd()))
	}

	return nil
}
