package scheduler

import (
	"context"
	"fmt"
	"spoti/internal/repository/postgres"

	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type EventWorker struct {
	repo      *postgres.Repository
	maxWorker int
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewEventWorker(ctx context.Context, repo *postgres.Repository) *EventWorker {
	ctx, cancel := context.WithCancel(ctx)
	return &EventWorker{
		repo:      repo,
		maxWorker: 100,
		ctx:       ctx,
		cancel:    cancel,
	}
}
func (ew *EventWorker) Start() {

}

func (ew *EventWorker) EventWorker() {
	log.Info().Msg("Starting task manager")

	ticker := time.NewTicker(15 * time.Minute)
	sem := make(chan struct{}, 36)

	var wg sync.WaitGroup

	for {
		select {
		case <-ew.ctx.Done():
			// close(sem)
			log.Info().Msg("Stopping processing")
			return
		case t := <-ticker.C:
			fmt.Println("Tick at", t)

			for i := range 10 {
				sem <- struct{}{}
				wg.Add(1)

				go func(i int) {
					defer func() {
						<-sem
						wg.Done()
					}()

					log.Info().Int("index", i).Msg("Stopping processing")
				}(i)
			}

		}

		wg.Wait()
	}

	time.Sleep(500 * time.Millisecond)
	ticker.Stop()
}

func (ew *EventWorker) Stop() {
	ew.cancel()
	time.Sleep(500 * time.Millisecond)
	log.Info().Msg("Stop EventWorker.")
}
