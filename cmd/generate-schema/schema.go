package main

import (
	"flag"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/m-lab/go/cloud/bqx"
	"github.com/m-lab/go/rtx"
	"github.com/m-lab/packet-test/api"
)

var (
	pair1Schema  string
	train1Schema string
)

func init() {
	flag.StringVar(&pair1Schema, "pair1", "/var/spool/datatypes/pair1.json", "Path to write the pair1 schema out to.")
	flag.StringVar(&train1Schema, "train1", "/var/spool/datatypes/train1.json", "Path to write the train1 schema out to.")
}

func main() {
	flag.Parse()

	pair1Result := api.Result{}
	schema, err := bigquery.InferSchema(pair1Result)
	rtx.Must(err, "Failed to generate pair1 schema")

	schema = bqx.RemoveRequired(schema)
	json, err := schema.ToJSONFields()
	rtx.Must(err, "Failed to marshal pair1 schema")

	err = os.WriteFile(pair1Schema, json, 0o644)
	rtx.Must(err, "Failed to write pair1 schema")

	train1Result := api.Result{}
	schema, err = bigquery.InferSchema(train1Result)
	rtx.Must(err, "Failed to generate train1 schema")

	schema = bqx.RemoveRequired(schema)
	json, err = schema.ToJSONFields()
	rtx.Must(err, "Failed to marshal train1 schema")

	err = os.WriteFile(train1Schema, json, 0o644)
	rtx.Must(err, "Failed to write train1 schema")
}
