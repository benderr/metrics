package storage

import (
	"context"
	"database/sql"

	"github.com/benderr/metrics/internal/server/config"
	"github.com/benderr/metrics/internal/server/dump"
	"github.com/benderr/metrics/internal/server/repository"
	"github.com/benderr/metrics/internal/server/repository/dbstorage"
	"github.com/benderr/metrics/internal/server/repository/filestorage"
	"github.com/benderr/metrics/internal/server/repository/inmemory"
)

func New(ctx context.Context, config *config.Config, logger repository.Logger) (repository.MetricRepository, func()) {
	var repo repository.MetricRepository
	closer := func() {}
	switch {
	case config.DatabaseDsn != "":
		db, dberr := sql.Open("pgx", config.DatabaseDsn)
		if dberr != nil {
			db.Close()
			panic(dberr)
		}

		dbRepo := dbstorage.NewWithRetry(db, logger)
		if err := dbRepo.Prepare(ctx); err != nil {
			db.Close()
			panic(err)
		}
		closer = func() {
			db.Close()
		}
		repo = dbRepo

	case config.FileStoragePath != "":
		sync := config.StoreInterval == 0
		fs := filestorage.New(config.FileStoragePath, sync, logger)
		dumper := dump.New(fs.Sync)
		if !sync {
			go dumper.Start(ctx, config.StoreInterval)
		}
		if config.Restore {
			fs.Restore(ctx)
		}

		repo = fs

	default:
		repo = inmemory.New()
	}
	return repo, closer
}
