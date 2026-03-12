package elastic

import (
	"context"
	"log"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/esdsl"
)

// type ElastiDB struct {
// 	client *elasticsearch.TypedClient
// }

func NewElasticDB(url string) *elasticsearch.TypedClient {
	client := elasticConn(url)
	if err := mappingIndex(client, spotiIndexName); err != nil {
		log.Panicf("err Mapping elastic: %v", err)
		return nil
	}

	return client
}

func elasticConn(url string) *elasticsearch.TypedClient {
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{
			url,
		},
	})
	if err != nil {
		log.Fatalf("err elastic connection %s", err)
	}
	return client
}

func mappingIndex(client *elasticsearch.TypedClient, indexName string) error {
	exists, err := client.Indices.Exists(indexName).IsSuccess(context.Background())
	if err != nil {
		return err
	}
	if exists {
		log.Println("index exists")
		return nil
	}

	mappings := esdsl.NewTypeMapping().
		AddProperty("id", esdsl.NewKeywordProperty().Index(false)).
		AddProperty("type", esdsl.NewKeywordProperty()).
		AddProperty("name", esdsl.NewSearchAsYouTypeProperty())

	_, err = client.Indices.Create(indexName).
		Mappings(mappings).
		Do(context.Background())
	if err != nil {
		return err
	}

	log.Println("index created")
	return nil
}
