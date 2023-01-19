package monitor

import (
	"context"
	"net/http"
	"time"

	"watchmen/repository"
)

type Monitor interface {
	CheckLinkHealth(ctx context.Context) error
}

type monitor struct {
	linkRepo repository.LinkRepository
	client   *http.Client
}

// NewMonitor creates a monitor instance
func NewMonitor(linkRepo repository.LinkRepository, timeout time.Duration) Monitor {
	mnt := new(monitor)
	mnt.linkRepo = linkRepo
	mnt.client = &http.Client{Timeout: timeout}

	return mnt
}

func (mnt *monitor) CheckLinkHealth(ctx context.Context) error {
	links, err := mnt.linkRepo.GetAllLinks(ctx)
	if err != nil {
		return err
	}

	for _, link := range links {
		res, err := mnt.client.Get(link.URL)
		if err != nil {
			return err
		}

		status := false
		if res.Status[0] == 50 {
			status = true
		}

		report := repository.LinkReport{
			LinkID: link.ID,
			Status: status,
		}

		err = mnt.linkRepo.CreateLinkReport(ctx, &report)
		if err != nil {
			return err
		}
	}

	return nil
}
