package prepare

import (
	"codebugs/log"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func watchSignal() {
	sc := make(chan os.Signal, 3) //nolint:gomnd
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGUSR2,
	)

	go func() {
		for {
			sig := <-sc
			handleSignal(sig)
		}
	}()
}

func handleSignal(sig os.Signal) {
	if sig == syscall.SIGUSR2 {
		log.Log().Info("Got signal to reopen log file", zap.Stringer("signal", sig))
	} else {
		log.Log().Info("Got signal to exit", zap.Stringer("signal", sig))

		err := log.Log().Sync()

		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Failed to sync logger", err)
			os.Exit(1)
		}

		os.Exit(0)
	}
}
