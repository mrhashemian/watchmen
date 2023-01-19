package monitor

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

// Run creates a ticker to execute mnt.CheckLinkHealth() every d duration
func Run(ctx context.Context, mnt Monitor, d time.Duration) {
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			err := mnt.CheckLinkHealth(ctx)
			if err != nil {
				log.Error(err)
			}
		case <-ctx.Done():
			log.Info("done signal received")
			ticker.Stop()
			return
		}
	}
}
