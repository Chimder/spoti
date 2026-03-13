package meilisearchrepo

import (
	"context"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
)

type MeiliRepository struct {
	client meilisearch.ServiceManager
}

func NewMeiliRepository(client meilisearch.ServiceManager) *MeiliRepository {
	return &MeiliRepository{client: client}
}

type Document struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type SearchResult = Document

func (r *MeiliRepository) Add(ctx context.Context, doc Document) error {
	index := r.client.Index(spotiIndexName)

	_, err := index.AddDocuments([]Document{doc}, nil)
	if err != nil {
		return fmt.Errorf("add doc: %w", err)
	}

	return nil
}

func (r *MeiliRepository) Search(ctx context.Context, query string) ([]SearchResult, error) {
	return r.search(ctx, query, "")
}

func (r *MeiliRepository) SearchByType(ctx context.Context, query, docType string) ([]SearchResult, error) {
	return r.search(ctx, query, docType)
}

func (r *MeiliRepository) search(ctx context.Context, query, docType string) ([]SearchResult, error) {
	index := r.client.Index(spotiIndexName)

	filter := ""
	if docType != "" {
		filter = fmt.Sprintf(`type = "%s"`, docType)
	}

	searchRes, err := index.SearchWithContext(ctx, query, &meilisearch.SearchRequest{
		Limit:  10,
		Filter: filter,
	})
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	results := make([]SearchResult, 0, len(searchRes.Hits))
	for _, hit := range searchRes.Hits {
		var doc SearchResult
		if err := hit.DecodeInto(&doc); err != nil {
			return nil, fmt.Errorf("decode hit: %w", err)
		}

		results = append(results, doc)
	}
	return results, nil
}
