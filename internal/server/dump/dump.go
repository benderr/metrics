package dump

import (
	"context"
	"time"
)

type SyncFunc func(ctx context.Context) error
type Dumper struct {
	sync SyncFunc
}

func New(sync SyncFunc) *Dumper {
	return &Dumper{
		sync: sync,
	}
}

func (d *Dumper) Start(ctx context.Context, storeIntervalSeconds int) {
	if storeIntervalSeconds == 0 {
		return
	}
	saveTicker := time.NewTicker(time.Second * time.Duration(storeIntervalSeconds))

	defer saveTicker.Stop()

	for {
		<-saveTicker.C
		go d.sync(ctx)
	}
}
