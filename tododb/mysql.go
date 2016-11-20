package tododb

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLDB struct {
	database   string
	password   string
	user       string
	appVersion string
}

func NewMySQLDB(config map[string]string, appVersion string) MySQLDB {
	if _, exists := config["database"]; !exists {
		config["database"] = "mysql"
	}

	if _, exists := config["database"]; !exists {
		config["password"] = "root"
	}

	if _, exists := config["database"]; !exists {
		config["user"] = "root"
	}

	return MySQLDB{
		database:   config["database"],
		password:   config["password"],
		user:       config["user"],
		appVersion: appVersion,
	}
}

const (
	mysqlDatabase   = "todo"
	mysqlTable      = "TodoTable"
	mysqlTodoColumn = "Todo"
)

var _ TodoDB = MySQLDB{}

func (mysqlDB MySQLDB) createMySQLClient() (*sql.DB, error) {
	return sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", mysqlDB.user, mysqlDB.password, mysqlDB.database, mysqlDatabase))
}

//TODO initalize -> check if table exists elxe create table WHERE

// SELECT table_name FROM information_schema.tables where table_schema='mysqlDB';
// if table in rows okay else create new table

func (mysqlDB MySQLDB) GetAllTodos() ([]string, error) {
	db, err := mysqlDB.createMySQLClient()
	if err != nil {
		return []string{}, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM %s", mysqlTodoColumn, mysqlTable))
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()
	todos := []string{}

	for rows.Next() {
		var todo string
		err := rows.Scan(&todo)
		if err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
	}

	if rows.Err() != nil {
		return todos, rows.Err()
	}

	return todos, nil
}

func (mysqlDB MySQLDB) SaveTodo(todo string) error {
	db, err := mysqlDB.createMySQLClient()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Query(fmt.Sprintf("INSERT INTO %s VALUES (%s)", mysqlTable, todo))
	return err
}

func (mysqlDB MySQLDB) DeleteTodo(todo string) error {
	db, err := mysqlDB.createMySQLClient()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Query(fmt.Sprintf("DELETE FROM %s WHERE %s=(%s)", mysqlTable, mysqlTodoColumn, todo))
	return err
}

func (mysqlDB MySQLDB) GetHealthStatus() map[string]string {
	//TODO implement
	return map[string]string{}
}

func (mysqlDB MySQLDB) RegisterMetrics() {
	//TODO implement
}

//todo create Test Case and container -> https://hub.docker.com/r/mysql/mysql-server/
// http://go-database-sql.org/accessing.html
