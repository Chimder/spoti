package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	UsersCount      = 300
	ArtistsCount    = 500
	AlbumsCount     = 1000
	RecordingsCount = 1000
	// TracksCount     = 1000
	PlaylistsCount = 150
)

var (
	userIDs   []uuid.UUID
	artistIDs []uuid.UUID
	albumIDs  []uuid.UUID
	// recordingIDs []uuid.UUID
	trackIDs    []uuid.UUID
	playlistIDs []uuid.UUID
)

func main() {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("err to conn to db: %v\n", err)
	}
	defer pool.Close()

	gofakeit.Seed(time.Now().UnixNano())

	now := time.Now()
	fmt.Println("run seed>>")
	if err := seedUsers(ctx, pool); err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n users OK")

	if err := seedArtists(ctx, pool); err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n artists OK")

	if err := seedAlbums(ctx, pool); err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n albums OK")

	if err := seedAlbumArtists(ctx, pool); err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n albumArtists OK")

	// if err := seedRecordings(ctx, pool); err != nil {
	// 	log.Fatal(err)
	// }

	if err := seedTracks(ctx, pool); err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n tracks OK")

	if err := seedArtistTracks(ctx, pool); err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n artistTracks OK")

	if err := seedUserSaveAlbums(ctx, pool); err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n us_save_alb OK")
	if err := seedPlaylists(ctx, pool); err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n playlist OK")

	if err := seedPlaylistTracks(ctx, pool); err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n playlist tracks OK")

	fmt.Printf("time %v\n", time.Since(now))
	printStatistics(ctx, pool)
}

func seedUsers(ctx context.Context, pool *pgxpool.Pool) error {
	for i := 0; i < UsersCount; i++ {
		user := GetFakeUser()
		var id uuid.UUID
		err := pool.QueryRow(ctx, `
			INSERT INTO users (user_name, email, image, followers, premium_status)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`,
			user.name, user.email, user.image, user.followers, user.premiumStatus,
		).Scan(&id)
		if err != nil {
			return fmt.Errorf("error seed user %d: %w", i, err)
		}
		userIDs = append(userIDs, id)
	}
	return nil
}

func seedArtists(ctx context.Context, pool *pgxpool.Pool) error {
	for i := range ArtistsCount {
		artist := GetFakeArtist()
		var id uuid.UUID
		err := pool.QueryRow(ctx, `
			INSERT INTO artists (url, uri, artist_name, image, followers, popularity, genres)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id
		`,
			artist.url, artist.uri, artist.artist_name, artist.image, artist.followers, artist.popularity, artist.genres,
		).Scan(&id)
		if err != nil {
			return fmt.Errorf("error seed artist %d: %w", i, err)
		}
		artistIDs = append(artistIDs, id)
	}
	return nil
}

func seedAlbums(ctx context.Context, pool *pgxpool.Pool) error {
	for i := range AlbumsCount {
		album := GetFakeAlbums()
		var id uuid.UUID

		err := pool.QueryRow(ctx, `
			INSERT INTO albums (album_type, total_tracks, image, album_name, uri, copyrights, album_label, popularity, release_date)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id
		`, album.album_label, album.total_tracks, album.image, album.album_name, album.uri, album.copyrights, album.album_label,
			album.popularity, album.release_date,
		).Scan(&id)
		if err != nil {
			return fmt.Errorf("error seeding album %d: %w", i, err)
		}
		albumIDs = append(albumIDs, id)
	}
	return nil
}

func seedAlbumArtists(ctx context.Context, pool *pgxpool.Pool) error {
	for _, albumID := range albumIDs {
		numArtists := gofakeit.Number(1, 3)
		selectedArtists := make(map[uuid.UUID]bool)

		for range numArtists {
			artistID := artistIDs[rand.Intn(len(artistIDs))]

			if selectedArtists[artistID] {
				continue
			}
			selectedArtists[artistID] = true

			_, err := pool.Exec(ctx, `
				INSERT INTO album_artists (album_id, artist_id)
				VALUES ($1, $2)
			`, albumID, artistID)
			if err != nil {
				return fmt.Errorf("err seed album_artists: %w", err)
			}
		}
	}
	return nil
}

func seedTracks(ctx context.Context, pool *pgxpool.Pool) error {
	for _, albumID := range albumIDs {

		var totalTracks int
		err := pool.QueryRow(ctx,
			`SELECT total_tracks FROM albums WHERE id = $1`,
			albumID,
		).Scan(&totalTracks)
		if err != nil {
			return err
		}

		numDiscs := gofakeit.Number(1, 2)

		trackNum := 1
		for disc := 1; disc <= numDiscs; disc++ {

			tracksInDisc := totalTracks / numDiscs
			if disc == numDiscs {
				tracksInDisc += totalTracks % numDiscs
			}

			for i := 0; i < tracksInDisc; i++ {

				isrc := fmt.Sprintf("US%s%02d%05d",
					gofakeit.LetterN(3),
					gofakeit.Number(0, 99),
					gofakeit.Number(0, 99999))

				var recordingID uuid.UUID
				err := pool.QueryRow(ctx, `
					INSERT INTO recordings (isrc, duration_ms, popularity, play_count, audio_uri)
					VALUES ($1, $2, 0, 0, $3)
					RETURNING id
				`,
					isrc,
					gofakeit.Number(120_000, 300_000),
					fmt.Sprintf("spotify:recording:%s", uuid.New()),
				).Scan(&recordingID)
				if err != nil {
					return err
				}

				var id uuid.UUID
				err = pool.QueryRow(ctx, `
						INSERT INTO tracks (
						album_id, recording_id, track_name, track_number, disc_number, explicit, is_playable, track_type, uri, islocal
						)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			      RETURNING id
					`, albumID, recordingID,
					fmt.Sprintf("%s %s", gofakeit.Adjective(), gofakeit.Noun()),
					trackNum, disc, gofakeit.Bool(), true, "track",
					fmt.Sprintf("spotify:track:%s", uuid.New()), false,
				).Scan(&id)
				if err != nil {
					return fmt.Errorf("err seed track  %w", err)
				}

				trackNum++
				trackIDs = append(trackIDs, id)
			}
		}
	}
	return nil
}

func seedArtistTracks(ctx context.Context, pool *pgxpool.Pool) error {
	rows, err := pool.Query(ctx, `SELECT id FROM tracks`)
	if err != nil {
		return err
	}
	defer rows.Close()

	batch := &pgx.Batch{}

	for rows.Next() {
		var trackID uuid.UUID
		if err := rows.Scan(&trackID); err != nil {
			return err
		}

		numArtists := gofakeit.Number(1, 3)
		selected := map[uuid.UUID]bool{}

		for range numArtists {
			artistID := artistIDs[rand.Intn(len(artistIDs))]
			if selected[artistID] {
				continue
			}
			selected[artistID] = true

			batch.Queue(`
				INSERT INTO artist_tracks (artist_id, track_id)
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING
			`, artistID, trackID)
		}
	}

	br := pool.SendBatch(ctx, batch)
	defer br.Close()

	return br.Close()
}

func seedUserSaveAlbums(ctx context.Context, pool *pgxpool.Pool) error {

	for _, userId := range userIDs {
		numAlbums := gofakeit.Number(1, 25)

		for range numAlbums {
			albumId := albumIDs[rand.Intn(len(albumIDs))]

			query := `
			INSERT INTO user_saved_albums (album_id, user_id) VALUES ($1, $2)
			ON CONFLICT DO NOTHING`
			_, err := pool.Exec(ctx, query, albumId, userId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func seedPlaylists(ctx context.Context, pool *pgxpool.Pool) error {
	for i := range PlaylistsCount {
		ownerID := userIDs[rand.Intn(len(userIDs))]

		var id uuid.UUID
		err := pool.QueryRow(ctx, `
			INSERT INTO playlists (owner_id, playlist_name, description, image, is_public, total)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`,
			ownerID,
			gofakeit.Sentence(2),
			gofakeit.Sentence(10),
			fmt.Sprintf("https://i.scdn.co/image/%s", gofakeit.UUID()),
			gofakeit.Bool(),
			gofakeit.Number(0, 50),
		).Scan(&id)
		if err != nil {
			return fmt.Errorf("err seed playlist %d: %w", i, err)
		}

		playlistIDs = append(playlistIDs, id)
	}
	return nil
}

func seedPlaylistTracks(ctx context.Context, pool *pgxpool.Pool) error {
	for _, playlistID := range playlistIDs {
		var total int
		err := pool.QueryRow(ctx, `SELECT total FROM playlists WHERE id = $1`, playlistID).Scan(&total)
		if err != nil {
			return fmt.Errorf("err get playlist: %w", err)
		}

		if total == 0 {
			continue
		}

		selectedTracks := make(map[uuid.UUID]bool)
		position := 1

		for len(selectedTracks) < total {
			trackID := trackIDs[rand.Intn(len(trackIDs))]

			if selectedTracks[trackID] {
				continue
			}
			selectedTracks[trackID] = true

			addedAt := gofakeit.DateRange(
				time.Now().AddDate(-1, 0, 0),
				time.Now(),
			)

			_, err := pool.Exec(ctx, `
				INSERT INTO playlist_tracks (playlist_id, track_id, track_position, added_at)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT DO NOTHING
			`, playlistID, trackID, position, addedAt)
			if err != nil {
				return fmt.Errorf("err seed playlist_tracks: %w", err)
			}
			position++
		}
	}
	return nil
}

func printStatistics(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n📊 Database Statistics:")

	stats := []struct {
		table string
		query string
	}{
		{"Users", "SELECT COUNT(*) FROM users"},
		{"Artists", "SELECT COUNT(*) FROM artists"},
		{"Albums", "SELECT COUNT(*) FROM albums"},
		{"Recordings", "SELECT COUNT(*) FROM recordings"},
		{"Tracks", "SELECT COUNT(*) FROM tracks"},
		{"Album-Artist relations", "SELECT COUNT(*) FROM album_artists"},
		{"Artist-Track relations", "SELECT COUNT(*) FROM artist_tracks"},
		{"user_saved_album", "SELECT COUNT(*) FROM user_saved_albums"},
		{"Playlists", "SELECT COUNT(*) FROM playlists"},
		{"Playlist-Track relations", "SELECT COUNT(*) FROM playlist_tracks"},
	}

	for _, stat := range stats {
		var count int
		pool.QueryRow(ctx, stat.query).Scan(&count)
		fmt.Printf("  • %s: %d\n", stat.table, count)
	}
}
