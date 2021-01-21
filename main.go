package main

import "C"
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type customer struct {
	ID   int     `json:"id"`
	Name string  `json:"name"`
	DOB  string  `json:"dob"`
	Addr address `json:"addr"`
}

type address struct {
	ID         int    `json:"id"`
	StreetName string `json:"streetName"`
	City       string `json:"city"`
	State      string `json:"state"`
	CusID      int    `json:"cus_id"`
}

// Router global, so that handlers can access the routers
var router = mux.NewRouter()
var db *sql.DB
var err error

func connectDatabase() {
	db, err = sql.Open("mysql", "root:password@/customer_service")
	fmt.Println("Database connected.")
	if err != nil {
		panic(err)
	}
}

func getCustomerAll(w http.ResponseWriter, r *http.Request) {
	// Return all customers.
	var res []customer
	//w.Header().Set("Content-Type", "application/json")
	query := `SELECT * FROM Customer INNER JOIN Address ON Customer.ID = Address.CusID ORDER BY Customer.ID, Address.ID`
	rows, err := db.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
	}
	// Iterate through all customers.
	for rows.Next() {
		var c customer
		if err := rows.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.StreetName, &c.Addr.City, &c.Addr.State, &c.Addr.CusID); err != nil {
			log.Fatal(err)
		}
		res = append(res, c)
	}
	_ = json.NewEncoder(w).Encode(res)
}

func getCustomerByID(w http.ResponseWriter, r *http.Request) {
	/*
		https://stackoverflow.com/questions/31622052/how-to-serve-up-a-json-response-using-go
		w.Header().Set("Content-Type", "application/json")
		Take all variables in the multiplexer as params.
		That allows us to pass that variable in JSON as ? in query param of SQL.
	*/
	params := mux.Vars(r)
	row := db.QueryRow("SELECT * FROM Customer INNER JOIN Address ON Customer.ID = Address.CusID and Customer.ID = ? ORDER BY Customer.ID, Address.ID;", params["id"])
	var c customer
	if err := row.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.StreetName, &c.Addr.City, &c.Addr.State, &c.Addr.CusID); err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// JSON data, combines marshal and writer.
	_ = json.NewEncoder(w).Encode(c)
}

func getCustomerByName(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	// Take all variables in the multiplexer as params.
	// That allows us to pass that variable in JSON as ? in query param of SQL.
	params := mux.Vars(r)
	rows, err := db.Query("SELECT * FROM Customer INNER JOIN Address ON Customer.ID = Address.CusID and Customer.Name = ? ORDER BY Customer.ID, Address.ID;", params["name"])
	if err != nil {
		panic(err.Error())
	}
	var res []customer // Changing tc
	for rows.Next() {
		var c customer
		if err := rows.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.StreetName, &c.Addr.City, &c.Addr.State, &c.Addr.CusID); err != nil {
			log.Fatal(err)
		}

		res = append(res, c)
	}
	_ = json.NewEncoder(w).Encode(res)

}

/*
func deleteCustomer, which deletes /customer/{id}. in the database belonging to an ID.
In the function, just return string "SUCCESS" in the implementation.
We dont want to delete the data.
*/
func NoContent(w http.ResponseWriter, r *http.Request) {
	// Set up any headers you want here.
	w.WriteHeader(http.StatusNoContent) // send the headers with a 204 response code.
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	// Setting up content tye for the response to JSON.
	w.Header().Set("Content-Type", "application/json")
	// Query parameters are route variables for the current request.
	params := mux.Vars(r)

	// Preparing a query to be executed. (May execute multiple statements here as well)
	// Get the data
	// getCustomerByID(w, r)
	// Delete not working on DB.
	stmt, err := db.Prepare("DELETE FROM Customer WHERE ID = ?;")
	if err != nil {
		panic(err.Error())
	}
	_, _ = stmt.Exec(params["id"])
	NoContent(w, r)
	//fmt.Fprintf(w, "Post with ID = %s was deleted", params["id"])
	//io.WriteString(os.Stdout, fmt.Sprintf("Post with ID = %s was deleted", params["id"]))
	// Internally converts given string to []byte/ JSON.
	//_, _ = io.WriteString(w, "SUCCESS")
}

//--------------------------------------------------------------------------------------------------
/*
// AgeAt gets the age of an entity at a certain time.
func AgeAt(birthDate time.Time, now time.Time) int {
	// Get the year number change since the player's birth.
	years := now.Year() - birthDate.Year()

	// If the date is before the date of birth, then not that many years have elapsed.
	birthDay := getAdjustedBirthDay(birthDate, now)
	if now.YearDay() < birthDay {
		years--
	}

	return years
}

// Age is shorthand for AgeAt(birthDate, time.Now()), and carries the same usage and limitations.
func Age(birthDate time.Time) int {
	return AgeAt(birthDate, time.Now())
}

// Gets the adjusted date of birth to work around leap year differences.
func getAdjustedBirthDay(birthDate time.Time, now time.Time) int {
	birthDay := birthDate.YearDay()
	currentDay := now.YearDay()
	if isLeap(birthDate) && !isLeap(now) && birthDay >= 60 {
		return birthDay - 1
	}
	if isLeap(now) && !isLeap(birthDate) && currentDay >= 60 {
		return birthDay + 1
	}
	return birthDay
}

// Works out if a time.Time is in a leap year.
func isLeap(date time.Time) bool {
	year := date.Year()
	if year%400 == 0 {
		return true
	} else if year%100 == 0 {
		return false
	} else if year%4 == 0 {
		return true
	}
	return false
}

func calculateAge(age string) int {
	// https://golangbyexample.com/get-age-given-dob-go/
	temp := strings.Split(age, "-")
	year, _ := strconv.Atoi(temp[2])
	month, _ := strconv.Atoi(temp[1])
	day, _ := strconv.Atoi(temp[0])
	dob := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return Age(dob)
}

func getDOB(year, month, day int) time.Time {
	dob := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return dob
}
*/

func dateInSeconds(d1 string) int {
	d1Slice := strings.Split(d1, "/")

	newDate := d1Slice[2] + "/" + d1Slice[1] + "/" + d1Slice[0]
	myDate, err := time.Parse("2006/01/02", newDate)

	if err != nil {
		panic(err)
	}

	return int(time.Now().Unix() - myDate.Unix())
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	var idr []interface{}
	var c customer
	query := `insert into Customer (Name,DOB) values(?,?)`
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &c)

	if c.Name == "" || c.DOB == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idr = append(idr, &c.Name)
	idr = append(idr, &c.DOB)
	age := dateInSeconds(c.DOB)
	if age/(365*24*3600) < 18 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "not eligible")
		return
	}
	row, er := db.Exec(query, idr...)
	if er != nil {

	}
	query = `insert into Address (StreetName, City, State, CusID) values(?,?,?,?)`
	var idd []interface{}
	if c.Addr.StreetName == "" || c.Addr.City == "" || c.Addr.State == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idd = append(idd, &c.Addr.StreetName)
	idd = append(idd, &c.Addr.City)
	idd = append(idd, &c.Addr.State)

	id, err := row.LastInsertId()
	idd = append(idd, id)

	_, err = db.Exec(query, idd...)
	if err != nil {

	}
	query = `select * from Customer inner join Address on Customer.ID = Address.CusID where Customer.ID=?`

	rows, err := db.Query(query, id)
	var cust customer
	for rows.Next() {

		rows.Scan(&cust.ID, &cust.Name, &cust.DOB, &cust.Addr.ID, &cust.Addr.StreetName, &cust.Addr.City, &cust.Addr.State, &cust.Addr.CusID)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cust)
}

func putCustomer(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var c customer
	err := json.Unmarshal(body, &c)
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	if len(c.DOB) == 0 {
		param := mux.Vars(r)
		id := param["id"]
		if c.Name != "" {
			_, err := db.Exec("update Customer set Name=? where ID=?", c.Name, id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
		}

		test := address{}
		if c.Addr != test {
			var data []interface{}
			query := "update Address set "
			if c.Addr.City != "" {
				query += "City = ? ,"
				data = append(data, c.Addr.City)
			}
			if c.Addr.State != "" {
				query += "State = ? ,"
				data = append(data, c.Addr.State)
			}
			if c.Addr.StreetName != "" {
				query += "StreetName = ? ,"
				data = append(data, c.Addr.StreetName)
			}
			query = query[:len(query)-1]
			query += "where CusID = ? and ID = ?"
			data = append(data, id)
			data = append(data, c.Addr.ID)
			_, err = db.Exec(query, data...)

		}
		json.NewEncoder(w).Encode(c)

	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	db, err = sql.Open("mysql", "root:password@/customer_service")
	fmt.Println("Database connected.")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Retrieve
	router.HandleFunc("/customer", getCustomerAll).Methods(http.MethodGet)
	router.HandleFunc("/customer/{id:[0-9]+}", getCustomerByID).Methods(http.MethodGet)
	router.HandleFunc("/customer/{name:[a-zA-Z ]+}", getCustomerByName).Methods(http.MethodGet)
	// Delete
	router.HandleFunc("/customer/{id:[0-9]+}", deleteCustomer).Methods(http.MethodDelete)
	// Create
	router.HandleFunc("/customer", createCustomer).Methods(http.MethodPost)
	// Update
	router.HandleFunc("/customer/{id:[0-9]+}", putCustomer).Methods(http.MethodPut)

	log.Fatal(http.ListenAndServe(":8000", router))
	/*
		if err := http.ListenAndServe(":8000", router); err != nil {
			// Handle error properly in your app.
			err1 := errors.New("problem spawning port")
			log.Fatal(err1, err)
		}

	*/

}
