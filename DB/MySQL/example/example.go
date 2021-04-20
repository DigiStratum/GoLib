package main

import (
	"fmt"

	db "github.com/DigiStratum/GoLib/DB"
	mysql "github.com/DigiStratum/GoLib/DB/MySQL"
)

func main() {
	//test_append_clone()
	test_mysql_query()
	//test_modify_passed_struct()
}

type todo struct {
	//Result
	Id	int
	Task	string
	Due	string
}

func Todo(id int, task string, due string) todo {
	t := todo{
		Id:	id,
		Task:	task,
		Due:	due,
	}
	return t
}

/*

-- MySQL Test database setup for "todolist"
create database todolist;
user todolist;
create table todo (
	id int primary key not null auto_increment,
	task varchar(250),
	due datetime
);
create user 'username'@'localhost' identified by 'password';
grant all on todolist.* to 'username'@'localhost';

*/


func test_modify_passed_struct() {
	t := Todo(1, "loaf", "asap")
	todos := modify_todo(t)
	for _, t := range todos {
		print_todo(t)
	}
}

func modify_todo(t todo) []todo {
	todos := []todo{}
	todos = append(todos, t)
	newId := 2
	t.Id = newId
	newDue := "nope"
	t.Due = newDue
	todos = append(todos, t)
	return todos
}

func print_todo(t todo) {
	fmt.Printf(
		"todo: { \"id\": \"%d\", \"task\": \"%s\", \"due\": \"%s\" }\n",
		t.Id, t.Task, t.Due,
	)
}

func test_mysql_query() {
	dsn := db.MakeDSN("username", "password", "localhost", "3306", "todolist")
	fmt.Printf("MySQL DSN is: %s\n", dsn)
	manager := mysql.NewManager()

	dbKey, err := manager.Connect(dsn)
	if nil != err {
		fmt.Printf("Error connecting: %s\n", err.Error())
	}

	dbConn, err := manager.GetConnection(*dbKey)
	if nil != err {
		fmt.Printf("Error getting connection: %s\n", err.Error())
	}
	query := mysql.NewQuery(
		"SELECT id, task, due FROM todo;",
		Todo(0,"",""),
	)

	results, err := query.Run(dbConn)
	if nil != err {
		fmt.Printf("Query Error: %s\n", err.Error())
	} else {
		for index, result := range *results {
			if todoResult, ok := result.(todo); ok {
				print_todo(todoResult)
			} else {
				fmt.Printf("Error converting result record to todo{%d}\n", index)
			}
		}
	}

	manager.Disconnect(*dbKey)
}

type record struct {
	value		int
}

// This test snippet demonstrates that an instance of a struct is cloned when it is appended() to a slice.
func test_append_clone() {
	// A proper struct
	r := record{
		value: 1,
	}

	collection := []record{}

	collection = append(collection, r)
	r.value++
	collection = append(collection, r)
	r.value++
	collection = append(collection, r)

	for index, record := range collection {
		fmt.Printf("collection[%d].record.value = %d\n", index, record.value)
	}
}

