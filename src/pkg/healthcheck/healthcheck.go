package healthcheck

import (
	"fmt"
	"net/http"
	"time"
)

const maxPollingAttempts = 30
const pollingInterval = 5

func PollPubSub(port string) error {
	for attempt := 1; attempt <= maxPollingAttempts; attempt++ {
		// send get request until 200 response
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s", port))
		if err == nil && resp.StatusCode == http.StatusOK {
			return nil
		}
		// Log the error and wait before the next polling attempt
		fmt.Printf("Attempt %d: Polling Pub/Sub client: %v\n", attempt, err)
		time.Sleep(pollingInterval * time.Second)
	}
	return fmt.Errorf("Pub/Sub service did not become available after %d attempts", maxPollingAttempts)
}
