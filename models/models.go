package models

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

type Person struct {
	name  string
	email string
}

type Journal struct {
	Name  string
	Entry string
}

type DataModel struct {
	DB *sql.DB
}

func (m *DataModel) Init() error {
	fmt.Printf("Opening database...\n")
	var err error
	m.DB, err = sql.Open("sqlite", "test.db")
	if err != nil {
		return err
	}

	if _, err = m.DB.Exec(`
CREATE TABLE IF NOT EXISTS person (
	user_id INTEGER PRIMARY KEY,
	name TEXT,
	email TEXT
);
CREATE TABLE IF NOT EXISTS journal (
	journal_id INTEGER PRIMARY KEY,
	user_id INTEGER,
	entry TEXT
);
CREATE VIRTUAL TABLE IF NOT EXISTS posts
USING FTS5(title, body);
`); err != nil {
		return err
	}
	return nil
}

func (m *DataModel) Close() { m.DB.Close() }

func (m *DataModel) LoadPersons() error {
	fmt.Printf("Adding persons...\n")
	if _, err := m.DB.Exec(`
INSERT INTO person (name, email)
VALUES
	('Bill', 'bdettelb@redhat.com'),
	('Tom', 'tsmith@redhat.com');
`); err != nil {
		return err
	}
	return nil
}

func (m *DataModel) LoadJournals() error {
	fmt.Printf("Adding journals...\n")
	if _, err := m.DB.Exec(`
INSERT INTO journal (user_id, entry)
VALUES
	(1, 'It was a dark and story night...'),
	(2, 'This is the time that tries mens souls');
`); err != nil {
		return err
	}
	return nil
}

func (m *DataModel) GetJournals() ([]Journal, error) {

	rows, err := m.DB.Query(`
SELECT name, entry
FROM journal
JOIN person ON journal.user_id = person.user_id;`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var journals []Journal

	for rows.Next() {
		var j Journal
		err := rows.Scan(&j.Name, &j.Entry)
		if err != nil {
			return nil, err
		}
		journals = append(journals, j)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return journals, nil
}

func (m *DataModel) LoadPostsFile() error {
	fmt.Printf("Adding war and peace...\n")
	readFile, err := os.Open("war_and_peace.txt")
	if err != nil {
		return err
	}
	defer readFile.Close()
	count := 0
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	txn, err := m.DB.Begin()
	if err != nil {
		return err
	}
	insertStmt := "INSERT INTO posts(title, body) VALUES (?, ?)"
	for fileScanner.Scan() {
		_, err = txn.Exec(insertStmt, "MyTitle", fileScanner.Text())
		if err != nil {
			txn.Rollback()
			return err
		}
		count++
	}
	if err := txn.Commit(); err != nil {
		return err
	}
	fmt.Printf("Loaded in %d lines\n", count)
	return nil
}

func (m *DataModel) SearchPosts(search string) error {

	rows, err := m.DB.Query("SELECT highlight(posts, 1, '<b>', '</b>') FROM posts WHERE posts MATCH '" + search + "';")

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var body string

		err := rows.Scan(&body)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", body)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}
