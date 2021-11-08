package graphql

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
)

const (
	defaultURL = "http://127.0.0.1:8080/graphql"
)

var sampleConfig = `
	url = <ip>:<port>/graphql
	query = <your_query>
`

type GraphQL struct {
	URL   string          `toml:"url"`
	Query string          `toml:"query"`
	Log   telegraf.Logger `toml:"-"`
}

func (g *GraphQL) Connect() error {
	return nil
}

func (g *GraphQL) Close() error {
	return nil
}

func (g *GraphQL) Description() string {
	return "A plugin that can load metrics with a graphql mutation"
}

func (g *GraphQL) SampleConfig() string {
	return sampleConfig
}

func (g *GraphQL) Write(metrics []telegraf.Metric) error {
	graphql_req := make(map[string]interface{})
	raw_data := make(map[string]interface{})
	data := make(map[string]interface{})
	for _, v := range metrics {
		data["uuid"] = v.Tags()["topic"]
		for _, f := range v.FieldList() {
			raw_data[f.Key] = f.Value
		}
		raw_data_json_bytes, err := json.Marshal(raw_data)
		if err != nil {
			return err
		}
		data["raw_json"] = string(raw_data_json_bytes)
		vars := make(map[string]interface{})
		vars["data"] = data
		graphql_req["variables"] = vars
		query := g.Query
		graphql_req["query"] = query
		json_str, err := json.Marshal(graphql_req)
		if err != nil {
			return err
		}
		_, err = http.Post(g.URL, "application/json", bytes.NewBuffer(json_str))
		if err != nil {
			return err
		}
		if err != nil {
			log.Fatalln(err)
		}
	}
	return nil
}

func init() {
	outputs.Add("graphql", func() telegraf.Output {
		return &GraphQL{
			URL: defaultURL,
		}
	})
}
