package postgres

import (
	"context"
	"fmt"
	"github.com/Chimder/spoti/internal/domain/album"
	"github.com/Chimder/spoti/internal/domain/artist"
	"github.com/Chimder/spoti/internal/domain/playlist"
	"github.com/Chimder/spoti/internal/domain/recording"
	"github.com/Chimder/spoti/internal/domain/track"
	"github.com/Chimder/spoti/internal/domain/user"
	albumrepo "github.com/Chimder/spoti/internal/repository/postgres/album"
	artistrepo "github.com/Chimder/spoti/internal/repository/postgres/artist"
	"github.com/Chimder/spoti/internal/repository/postgres/pgiface"
	playlistrepo "github.com/Chimder/spoti/internal/repository/postgres/playlist"
	recordingrepo "github.com/Chimder/spoti/internal/repository/postgres/recording"
	trackrepo "github.com/Chimder/spoti/internal/repository/postgres/track"
	userrepo "github.com/Chimder/spoti/internal/repository/postgres/user"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Repository struct {
	Db        pgiface.Querier
	Pool      *pgxpool.Pool
	User      user.UserRepository
	Artist    artist.ArtistRepository
	Album     album.AlbumRepository
	Playlist  playlist.PlaylistRepository
	Track     track.TrackRepository
	Recording recording.RecordingRepository
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return newRepository(pool, pool)
}

func newRepository(db pgiface.Querier, pool *pgxpool.Pool) *Repository {
	return &Repository{
		Db:        db,
		Pool:      pool,
		User:      userrepo.NewUserRepo(db),
		Artist:    artistrepo.NewArtistRepo(db),
		Album:     albumrepo.NewAlbumRepo(db),
		Playlist:  playlistrepo.NewPlaylistRepo(db),
		Recording: recordingrepo.NewRecordingRepo(db),
		Track:     trackrepo.NewTrackRepo(db),
	}
}

func (r *Repository) newWithTx(tx pgx.Tx) *Repository {
	return newRepository(tx, r.Pool)
}

func (r *Repository) WithTx(ctx context.Context, fn func(*Repository) error) (err error) {
	return r.WithTxOptions(ctx, pgx.TxOptions{}, fn)
}
func (r *Repository) WithTxOptions(ctx context.Context, opts pgx.TxOptions, fn func(*Repository) error) (err error) {
	tx, err := r.Pool.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("start tx: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Error().Err(rbErr).Msg("rollback error in recover")
			}
			log.Error().Any("panic", p).Msg("tx panic recover")
			panic(p)
		} else if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Error().Err(rbErr).Msg("rollback error")
			}
			log.Error().Err(err).Msg("tx rolled back")
		}
	}()

	txRepo := r.newWithTx(tx)

	if err = fn(txRepo); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Error().Err(err).Msg("commit failed")
		return fmt.Errorf("commit tx err: %w", err)
	}

	return nil
}
