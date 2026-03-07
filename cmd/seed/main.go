package main

import "spoti/internal/repository/clickhouse"

func main() {
	// postgres.RunSeed()
	clickhouse.StartSeedListeningEvents()
}
