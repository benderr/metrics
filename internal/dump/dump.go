package dump

import (
	"time"

	"github.com/benderr/metrics/internal/repository/filestorage"
)

type Dumper struct {
	fileRepo *filestorage.FileMetricRepository
}

func New(fileRepo *filestorage.FileMetricRepository) *Dumper {
	return &Dumper{
		fileRepo: fileRepo,
	}
}

func (d *Dumper) Start(storeIntervalSeconds int) {
	if storeIntervalSeconds == 0 {
		return
	}
	saveTicker := time.NewTicker(time.Second * time.Duration(storeIntervalSeconds))

	defer saveTicker.Stop()

	for {
		<-saveTicker.C
		go d.fileRepo.Sync()
	}
}
