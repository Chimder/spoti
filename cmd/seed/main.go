package main

import (
	"spoti/internal/repository/clickhouse"
	meilisearchrepo "spoti/internal/repository/meilisearch"
	"spoti/internal/repository/postgres"
)

func main() {
	postgres.RunSeed()
	clickhouse.StartSeedListeningEvents()
	meilisearchrepo.StartSeedMeiliSearch()
}
