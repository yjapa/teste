package templates

// Transactions sets the configuration for transaction index

var Block = Object{
	"index_patterns": Array{
		"blocks-*",
	},
	"settings": Object{
		"number_of_shards":   5,
		"number_of_replicas": 0,
		"index": Object{
			"sort.field": Array{
				"timestamp", "txId", "blockHeight",
			},
			"sort.order": Array{
				"desc", "desc",
			},
		},
	},
	"mappings": Object{
		"properties": Object{
			"hash": Object{
				"type": "string",
			},
			"blockHeight": Object{
				"type": "string",
			},
			"blockRewards": Object{
				"type": "string",
			},
			"size": Object{
				"type": "string",
			},
			"timestamp": Object{
				"type": "date",
			},
			"transactions": Object{
				"type": "nested",
				"properties": Object{
					"id": Object{
						"type": "long",
					},
					"txId": Object{
						"type": "string",
					},
					"block": Object{
						"type": "string",
					},
					"fromAddress": Object{
						"type": "string",
					},
					"toAddress": Object{
						"type": "string",
					},
					"value": Object{
						"type": "string",
					},
					"fee": Object{
						"type": "string",
					},
					"status": Object{
						"type": "integer",
					},
					"timestamp": Object{
						"type": "date",
					},
					"interactions": Object{
						"type": "nested",
						"properties": Object{
							"txId": Object{
								"type": "string",
							},
							"fromAddress": Object{
								"type": "string",
							},
							"toAddress": Object{
								"type": "string",
							},
							"value": Object{
								"type": "string",
							},
							"contractType": Object{
								"type": "string",
							},
							"contractName": Object{
								"type": "string",
							},
						},
					},
				},
			},
		},
	},
}
