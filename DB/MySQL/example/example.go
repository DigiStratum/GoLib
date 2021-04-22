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

type Todo struct {
	//Result
	Id	int
	Task	string
	Due	string
}

// Satisfies ResultIfc
func (t Todo) ZeroClone() (mysql.ResultIfc, mysql.PropertyPointers) {
	n := Todo{}
	npp := mysql.PropertyPointers{ &n.Id, &n.Task, &n.Due }
	return &n, npp
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
	t := Todo{ Id: 1, Task: "loaf", Due: "asap" }
	todos := modify_todo(t)
	for _, t := range todos {
		print_todo(t)
	}
}

func modify_todo(t Todo) []Todo {
	todos := []Todo{}
	todos = append(todos, t)
	t.Id = 2
	t.Due = "nope"
	todos = append(todos, t)
	return todos
}

func print_todo(t Todo) {
	fmt.Printf(
		"Todo: { \"id\": \"%d\", \"task\": \"%s\", \"due\": \"%s\" }\n",
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
		Todo{},
	)

	results, err := query.Run(dbConn)
	if nil != err {
		fmt.Printf("Query Error: %s\n", err.Error())
	} else {
		for index, result := range *results {
			if todoResult, ok := result.(*Todo); ok {
				print_todo(*todoResult)
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

