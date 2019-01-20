package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
func test() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/wordpress")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Println("here")
	prepareQuery(db)
}

func database(db *sql.DB) {
	var (
		id         int
		post_title string
	)
	rows, err := db.Query("select id, post_title from wp_posts where id = ?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &post_title)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, post_title)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func prepareQuery(db *sql.DB) {
	var (
		id         int
		post_title string
	)
	stmt, err := db.Prepare("select id, post_title from wp_posts where id >= ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &post_title)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, post_title)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}
