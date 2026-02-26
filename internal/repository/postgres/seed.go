package postgres

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"spoti/internal/repository/clickhouse"
	"spoti/internal/repository/postgres/testhelpers"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	UsersCount     = 300
	ArtistsCount   = 500
	AlbumsCount    = 1000
	PlaylistsCount = 150
)

var (
	userIDs     []uuid.UUID
	artistIDs   []uuid.UUID
	albumIDs    []uuid.UUID
	trackIDs    []uuid.UUID
	playlistIDs []uuid.UUID
)

func RunSeed() {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	repo := NewRepository(pool)
	gofakeit.Seed(time.Now().UnixNano())

	fmt.Println("Starting seed...")

	if err := seedUsers(ctx, repo); err != nil {
		log.Fatal(err)
	}
	if err := seedArtists(ctx, repo); err != nil {
		log.Fatal(err)
	}
	if err := seedAlbums(ctx, repo); err != nil {
		log.Fatal(err)
	}
	if err := seedAlbumArtists(ctx, repo); err != nil {
		log.Fatal(err)
	}
	if err := seedTracks(ctx, repo); err != nil {
		log.Fatal(err)
	}
	if err := seedArtistTracks(ctx, repo); err != nil {
		log.Fatal(err)
	}
	if err := seedUserSaveAlbums(ctx, repo); err != nil {
		log.Fatal(err)
	}
	if err := seedUserSavedArtists(ctx, repo); err != nil {
		log.Fatal(err)
	}
	if err := seedPlaylists(ctx, repo); err != nil {
		log.Fatal(err)
	}
	if err := seedPlaylistTracks(ctx, repo); err != nil {
		log.Fatal(err)
	}
	if err := seedUserSavedPlaylists(ctx, repo); err != nil {
		log.Fatal(err)
	}

	if err = clickhouse.SeedListeningEvents(ctx, repo.pool); err != nil {
		fmt.Println("Err seed clickhouse")
		log.Fatal(err)
	}

	printStatistics(ctx, repo.pool)
	fmt.Println("Seed success")
}

// Seed //

func seedUsers(ctx context.Context, repo *Repository) error {
	for range UsersCount {
		id, err := testhelpers.CreateUser(repo.User)
		if err != nil {
			return fmt.Errorf("seed user: %w", err)
		}
		userIDs = append(userIDs, id)
	}
	fmt.Println("Users seeded")
	return nil
}

func seedArtists(ctx context.Context, repo *Repository) error {
	for range ArtistsCount {
		id, err := testhelpers.CreateArtist(repo.Artist)
		if err != nil {
			return fmt.Errorf("seed artist: %w", err)
		}
		artistIDs = append(artistIDs, id)
	}
	fmt.Println("Artists seeded")
	return nil
}

func seedAlbums(ctx context.Context, repo *Repository) error {
	for range AlbumsCount {
		id, err := testhelpers.CreateAlbum(repo.Album)
		if err != nil {
			return fmt.Errorf("seed album: %w", err)
		}
		albumIDs = append(albumIDs, id)
	}
	fmt.Println("Albums seeded")
	return nil
}

func seedAlbumArtists(ctx context.Context, repo *Repository) error {
	for _, albumID := range albumIDs {
		numArtists := rand.Intn(3) + 1
		selected := make(map[uuid.UUID]bool)
		for range numArtists {
			artistID := artistIDs[rand.Intn(len(artistIDs))]
			if selected[artistID] {
				continue
			}
			selected[artistID] = true

			if err := repo.Album.AddArtistToAlbum(ctx, albumID, artistID); err != nil {
				return fmt.Errorf("seed album artists: %w", err)
			}
		}
	}
	fmt.Println("Album-Artists relations seeded")
	return nil
}

func seedTracks(ctx context.Context, repo *Repository) error {
	for _, albumID := range albumIDs {
		album, err := repo.Album.GetAlbum(ctx, albumID.String())
		if err != nil {
			return err
		}

		totalTracks := album.TotalTracks
		numDiscs := rand.Intn(2) + 1
		trackNum := 1

		for disc := 1; disc <= numDiscs; disc++ {
			tracksInDisc := totalTracks / numDiscs
			if disc == numDiscs {
				tracksInDisc += totalTracks % numDiscs
			}

			for i := 0; i < tracksInDisc; i++ {
				recID, err := testhelpers.CreateRecording(repo.Recording)
				if err != nil {
					return fmt.Errorf("seed recording: %w", err)
				}
				trackID, err := testhelpers.CreateTrack(repo.Track, albumID, recID, trackNum, disc)
				if err != nil {
					return fmt.Errorf("seed track: %w", err)
				}
				trackIDs = append(trackIDs, trackID)
				trackNum++
			}
		}
	}
	fmt.Println("Tracks seeded")
	return nil
}

func seedArtistTracks(ctx context.Context, repo *Repository) error {
	for _, trackID := range trackIDs {
		numArtists := rand.Intn(3) + 1
		selected := map[uuid.UUID]bool{}
		for range numArtists {
			artistID := artistIDs[rand.Intn(len(artistIDs))]
			if selected[artistID] {
				continue
			}
			selected[artistID] = true

			if err := repo.Track.AddArtistToTrack(ctx, trackID, artistID); err != nil {
				return fmt.Errorf("seed artist tracks: %w", err)
			}
		}
	}
	fmt.Println("Artist-Tracks relations seeded")
	return nil
}
func seedUserSaveAlbums(ctx context.Context, repo *Repository) error {

	for _, userId := range userIDs {
		numAlbums := gofakeit.Number(1, 25)

		for range numAlbums {
			albumId := albumIDs[rand.Intn(len(albumIDs))]

			query := `
			INSERT INTO user_saved_albums (album_id, user_id) VALUES ($1, $2)
			ON CONFLICT DO NOTHING`
			_, err := repo.pool.Exec(ctx, query, albumId, userId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func seedUserSavedArtists(ctx context.Context, repo *Repository) error {
	for _, userID := range userIDs {
		numArtists := gofakeit.Number(0, 15)

		if numArtists == 0 {
			continue
		}

		selectedArtists := make(map[uuid.UUID]bool)

		for range numArtists {
			artistID := artistIDs[rand.Intn(len(artistIDs))]

			if selectedArtists[artistID] {
				continue
			}
			selectedArtists[artistID] = true

			_, err := repo.pool.Exec(ctx, `
				INSERT INTO user_saved_artists (user_id, artist_id)
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING
			`, userID, artistID)
			if err != nil {
				return fmt.Errorf("err seed user_saved_artists: %w", err)
			}
		}
	}
	return nil
}

func seedUserSavedAlbums(ctx context.Context, repo *Repository) error {
	for _, userID := range userIDs {
		numAlbums := rand.Intn(25) + 1
		selected := map[uuid.UUID]bool{}
		for range numAlbums {
			albumID := albumIDs[rand.Intn(len(albumIDs))]
			if selected[albumID] {
				continue
			}
			selected[albumID] = true

			if err := repo.Album.SaveAlbumsForCurrentUser(ctx, []string{albumID.String()}, userID.String()); err != nil {
				return fmt.Errorf("seed user saved albums: %w", err)
			}
		}
	}
	fmt.Println("User saved albums seeded")
	return nil
}

func SeedUserSavedArtists(ctx context.Context, repo *Repository) error {
	for _, userID := range userIDs {
		numArtists := rand.Intn(16)
		selected := map[uuid.UUID]bool{}
		for range numArtists {
			artistID := artistIDs[rand.Intn(len(artistIDs))]
			if selected[artistID] {
				continue
			}
			selected[artistID] = true

			if err := repo.User.FollowUserToArtist(ctx, userID, artistID); err != nil {
				return fmt.Errorf("seed user saved artists: %w", err)
			}
		}
	}
	fmt.Println("User saved artists seeded")
	return nil
}

func seedPlaylists(ctx context.Context, repo *Repository) error {
	for range PlaylistsCount {
		ownerID := userIDs[rand.Intn(len(userIDs))]
		playlistID, err := testhelpers.CreatePlaylist(repo.Playlist, ownerID)
		if err != nil {
			return fmt.Errorf("seed playlist: %w", err)
		}
		playlistIDs = append(playlistIDs, playlistID)
	}
	fmt.Println("Playlists seeded")
	return nil
}

func seedPlaylistTracks(ctx context.Context, repo *Repository) error {
	for _, playlistID := range playlistIDs {
		total := rand.Intn(50) + 1
		selected := map[uuid.UUID]bool{}
		position := 1

		for len(selected) < total {
			trackID := trackIDs[rand.Intn(len(trackIDs))]
			if selected[trackID] {
				continue
			}
			selected[trackID] = true

			if err := repo.Playlist.AddToPlaylist(ctx, playlistID.String(), trackID.String()); err != nil {
				return fmt.Errorf("seed playlist tracks: %w", err)
			}
			position++
		}
	}
	fmt.Println("Playlist-Tracks seeded")
	return nil
}

func seedUserSavedPlaylists(ctx context.Context, repo *Repository) error {
	for _, userID := range userIDs {
		numPlaylists := rand.Intn(16)
		selected := map[uuid.UUID]bool{}
		for range numPlaylists {
			playlistID := playlistIDs[rand.Intn(len(playlistIDs))]
			if selected[playlistID] {
				continue
			}
			selected[playlistID] = true

			if err := repo.User.FollowUserToPlaylist(ctx, userID, playlistID); err != nil {
				return fmt.Errorf("seed user saved playlists: %w", err)
			}
		}
	}
	fmt.Println("User saved playlists seeded")
	return nil
}
func printStatistics(ctx context.Context, pool *pgxpool.Pool) {
	stats := []struct {
		table string
		query string
	}{
		{"users", "SELECT COUNT(*) FROM users"},
		{"artists", "SELECT COUNT(*) FROM artists"},
		{"albums", "SELECT COUNT(*) FROM albums"},
		{"recordings", "SELECT COUNT(*) FROM recordings"},
		{"tracks", "SELECT COUNT(*) FROM tracks"},
		{"album artist relations", "SELECT COUNT(*) FROM album_artists"},
		{"artist track relations", "SELECT COUNT(*) FROM artist_tracks"},
		{"user saved album", "SELECT COUNT(*) FROM user_saved_albums"},
		{"user saved artists", "SELECT COUNT(*) FROM user_saved_artists"},
		{"playlists", "SELECT COUNT(*) FROM playlists"},
		{"playlist track relations", "SELECT COUNT(*) FROM playlist_tracks"},
		{"user saved playlists", "SELECT COUNT(*) FROM user_saved_playlists"},
	}

	for _, stat := range stats {
		var count int
		pool.QueryRow(ctx, stat.query).Scan(&count)
		fmt.Printf("  • %s: %d\n", stat.table, count)
	}
}
