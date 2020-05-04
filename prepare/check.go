package prepare

import (
	"codebugs/config"
	"codebugs/log"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func PingServer() error {
	for i := 0; i < 10; i++ {
		// Ping the server by sending a GET request to `/health`.
		pingURL := "127.0.0.1:" + strconv.Itoa(config.Config.Serve.Port)
		log.Log().Info(pingURL)
		url := fmt.Sprintf("%s/health_check", pingURL)
		resp, err := http.Get(url) //nolint:gosec
		if err != nil {
			return err
		}
		if resp.StatusCode == http.StatusOK {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Log().Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("cannot connect to the router")
}

// check api run
func APICheck() {
	if err := PingServer(); err != nil {
		log.Log().Fatal("The router has no response, or it might took too long to start up.")
	}
	log.Log().Info("The router has been deployed successfully.")
}
