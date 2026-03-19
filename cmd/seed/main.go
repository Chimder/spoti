package main

import (
	"github.com/Chimder/spoti/internal/repository/clickhouse"
	meilisearchrepo "github.com/Chimder/spoti/internal/repository/meilisearch"
	"github.com/Chimder/spoti/internal/repository/postgres"
)

func main() {
	postgres.RunSeed()
	clickhouse.StartSeedListeningEvents()
	meilisearchrepo.StartSeedMeiliSearch()
}
