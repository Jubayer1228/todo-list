package middleware

/*
	This is the main handler file which will be responsible to create the database connection, CRUD operation etc
	createConnection function will read the .env file to get the postgres database connection info and for simplicity sslmode is disabled. db.Ping is the method that will create the connection whereas sql.Open() methos is just testing the connection  info 
	
*/
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"todo-list/models"
	"log"
	"net/http"
	"os"
	"path"
	"text/template"
	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // golang driver for postgres
)

type response struct {
	ID      string `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// create connection with postgres db
// createConnection function will read the .env file and generate a dbURL file and use that as param for sql.Open() methos which takes two arguments. One is driver type, in this case postgres is the type and the connection info as the other param. once the connection is verified it will create the connection
func createConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file.")
	}
	
	username, password, database, HOST, PORT :=
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username, password, HOST, PORT, database)
	
	// Open the connection
	db, err := sql.Open("postgres", dbURL)
	//db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Database is unsuccessful to connect")
		panic(err)
	}

	fmt.Println("Database Connected Successfully.")

	// return the connection
	return db
}

// HomePage Handler
// template package is used to parse the filepath
func HomePage(w http.ResponseWriter, r *http.Request) {
	var filepath = path.Join("views", "index.html")
	var tmpl, err = template.ParseFiles(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetAllTodoList will return all the todoList
func GetAllTodoList(w http.ResponseWriter, r *http.Request) {

	// get all the todos in the db
	todos, err := getAllTodosList()

	if err != nil {
		log.Fatalf("Unable to get all to-do list. %v", err)
	}

	// send all the todos as response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// CreateTodo create a todo list
func CreateTodoList(w http.ResponseWriter, r *http.Request) {
	// create an empty todo of
	var todo models.TodoList

	// decode the json request to todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	// call insert todo function and pass the todo
	insertID := insertTodoList(todo)

	// format a response object
	res := response{
		ID:      insertID,
		Message: "To-do list created successfully.",
	}

	// send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// UpdateTodo update todo's detail
func UpdateTodoList(w http.ResponseWriter, r *http.Request) {
	// create an empty todo of type models.TodoList
	var todo models.TodoList

	// decode the json request to todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call update todo list to update the list
	// 
	updatedRows := updateTodoList(todo)

	// format the message string 
	msg := fmt.Sprintf("User updated successfully. Total rows/record affected %v", updatedRows)

	// format the response message
	res := response{
		ID:      todo.ID,
		Message: msg,
	}

	// send the response as json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// DeleteTodo delete todo's detail
func DeleteTodoList(w http.ResponseWriter, r *http.Request) {

	// get the todo list id from the request params, key is "id"
	params := mux.Vars(r)
	id := params["id"]

	// call the deleteTodo
	deletedRows := deleteTodoList(id)

	// format the message string
	msg := fmt.Sprintf("Todo list deleted successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := response{
		ID:      id,
		Message: msg,
	}

	// send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// These are the  handler functions that will interact with the postgres database
// first, it will create the connection by calling the createConnection function
// once connection is verified, create a todos variable to hold the todos list by calling the model.TodoList structure which is a custom datatype for this TodoList
// then from the todos table get all the todosList, how to create a table is added on the sql script
// db.Query function takes any sql statement and it return a list and then we need to iterate the rows 
// rows.Scan function will take all the custom datatype arguments as input and append them as todos. 
// get list todo from the DB
func getAllTodosList() ([]models.TodoList, error) {
	// create db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var todos []models.TodoList

	// create the select sql query
	sqlStatement := `SELECT * FROM todos`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows and create the todos
	for rows.Next() {
		var todo models.TodoList

		// unmarshal the row object to todo
		err = rows.Scan(&todo.ID, &todo.Text, &todo.Checked)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the todo in the todos slice
		todos = append(todos, todo)

	}

	// return empty todo on error
	return todos, err
}

// insert one todo list in the DB
// 
func insertTodoList(todo models.TodoList) string {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning todoid will return the id of the inserted todo
	sqlStatement := `INSERT INTO todos (text, checked) VALUES ($1, $2) RETURNING id`

	// the inserted id will store in this id
	var id string

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := db.QueryRow(sqlStatement, todo.Text, todo.Checked).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

// update todo in the DB
func updateTodoList(todo models.TodoList) int64 {

	// create the db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE todos SET checked=$2 WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, todo.ID, todo.Checked)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// delete todo in the DB
func deleteTodoList(id string) int64 {

	// create the db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM todos WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
