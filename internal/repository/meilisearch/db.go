package meilisearchrepo

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/meilisearch/meilisearch-go"
	"github.com/rs/zerolog/log"
)

const spotiIndexName = "spoti"

func NewMeiliDB(ctx context.Context,url string) meilisearch.ServiceManager {
	client := meiliConn(url)
	if err := mappingIndex(ctx, client, spotiIndexName); err != nil {
		log.Panic().Str("err setup meilisearch", err.Error()).Msg("err mapping meili")
	}
	return client
}

func meiliConn(url string) meilisearch.ServiceManager {
	return meilisearch.New(
		url,
		meilisearch.WithCustomJsonMarshaler(sonic.Marshal),
		meilisearch.WithCustomJsonUnmarshaler(sonic.Unmarshal),
	)
}

func mappingIndex(ctx context.Context, client meilisearch.ServiceManager, indexName string) error {
	_, err := client.CreateIndex(ctx, &meilisearch.IndexConfig{
		Uid:        indexName,
		PrimaryKey: "id",
	})
	if err != nil {
		if meiliErr, ok := err.(*meilisearch.Error); ok && meiliErr.StatusCode == 409 {
			log.Error().Err(err).Int("code", meiliErr.StatusCode).Msg("err create meilisearch index")
		} else {
			return fmt.Errorf("create index: %w", err)
		}
	}

	index := client.Index(indexName)

	settings := &meilisearch.Settings{
		SearchableAttributes: []string{"name"},
		FilterableAttributes: []string{"type"},
	}

	if _, err := index.UpdateSettings(ctx, settings); err != nil {
		return fmt.Errorf("update settings: %w", err)
	}

	log.Info().Msg("meilisearch index ok")
	return nil
}
