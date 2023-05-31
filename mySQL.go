package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type MySQLData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	// Koneksi MySQL
	mysqlDB, err := sql.Open("mysql", "root:@tcp(localhost:3306)/mahasiswa")
	if err != nil {
		log.Fatal(err)
	}
	defer mysqlDB.Close()

	// Membuat schema GraphQL
	fields := graphql.Fields{
		"data": &graphql.Field{
			Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
				Name: "MySQLData",
				Fields: graphql.Fields{
					"id":   &graphql.Field{Type: graphql.Int},
					"name": &graphql.Field{Type: graphql.String},
				},
			})),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var data []MySQLData

				rows, err := mysqlDB.Query("SELECT * FROM person")
				if err != nil {
					return nil, err
				}
				defer rows.Close()

				for rows.Next() {
					var d MySQLData
					err := rows.Scan(&d.ID, &d.Name)
					if err != nil {
						return nil, err
					}
					data = append(data, d)
				}

				return data, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Handler GraphQL
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Mengatur rute GraphQL
	http.Handle("/graphql", h)

	// Menjalankan server
	fmt.Println("Server GraphQL berjalan di http://localhost:8080/graphql")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
