package elastic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/esdsl"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/textquerytype"
)

type ElasticRepository struct {
	client *elasticsearch.TypedClient
}

const spotiIndexName = "spoti"

func NewElasticRepository(client *elasticsearch.TypedClient) *ElasticRepository {
	return &ElasticRepository{client: client}
}

type Document struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type SearchResult = Document

func (r *ElasticRepository) Index(ctx context.Context, doc Document) error {
	resp, err := r.client.Index(spotiIndexName).
		Request(&doc).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("index document: %w", err)
	}

	if resp.Result.Name != "created" && resp.Result.Name != "updated" {
		return fmt.Errorf("Err create els index: %s", resp.Result.Name)
	}

	return nil
}

func (r *ElasticRepository) Search(ctx context.Context, query string) ([]SearchResult, error) {
	return r.search(ctx, query, "")
}

func (r *ElasticRepository) SearchByType(ctx context.Context, query, docType string) ([]SearchResult, error) {
	return r.search(ctx, query, docType)
}

func (r *ElasticRepository) search(ctx context.Context, query, docType string) ([]SearchResult, error) {
	searchQuery := esdsl.NewBoolQuery().
		Must(
			esdsl.NewMultiMatchQuery(query).
				Type(textquerytype.Boolprefix).
				Fields("name", "name._2gram", "name._3gram"),
		)

	if docType != "" {
		searchQuery = searchQuery.Filter(
			esdsl.NewTermQuery("type", esdsl.NewFieldValue().String(docType)),
		)
	}

	resp, err := r.client.Search().Index(spotiIndexName).Query(searchQuery).Size(10).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	results := make([]SearchResult, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var doc SearchResult
		if err := json.Unmarshal(hit.Source_, &doc); err != nil {
			return nil, fmt.Errorf("err unmarshal hit: %w", err)
		}
		results = append(results, doc)
	}
	return results, nil
}
