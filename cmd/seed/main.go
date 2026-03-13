package main

import (
	meilisearchrepo "spoti/internal/repository/meilisearch"
)

func main() {
	// postgres.RunSeed()
	// clickhouse.StartSeedListeningEvents()
	meilisearchrepo.StartSeedMeiliSearch()
}
