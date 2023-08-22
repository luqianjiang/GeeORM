package day1_database_sql

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"testing"
)

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "gee.db")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	if err == nil {
		affected, _ := result.RowsAffected()
		fmt.Println(affected)
	}
	row := engine.db.QueryRow("SELECT Name FROM User LIMIT 1")
	var name string
	if err := row.Scan(&name); err == nil {
		log.Println(name)
	}

	return engine
}

func TestNewEngine(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
}
