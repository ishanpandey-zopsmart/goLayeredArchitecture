package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetAllCustomers(t *testing.T) {
	var testCases = []struct {
		inp string
		out []customer
	}{
		{"", []customer{
			customer{1, "Ishan Pandey", "25-12-1999", address{1, "SSV Colony", "Lucknow", "UP", 1}},
			customer{2, "Varun Singh", "23-12-1996", address{2, "St.Peters Street", "Singapore", "Malaysia", 2}}},
		},
	}

	for i := range testCases {
		req := httptest.NewRequest("GET", "/customer/", nil)
		w := httptest.NewRecorder()
		getCustomerAll(w, req)

		resp := w.Body.Bytes()
		var res []customer
		_ = json.Unmarshal(resp, &res)
		if w.Code != http.StatusOK {
			t.Logf("FAILED: ERROR IN CONNECTION")
		} else if !reflect.DeepEqual(res, testCases[i].out) {
			t.Errorf("FAILED for %v expected %v got %v\n", testCases[i].inp, testCases[i].out, res)
		}

	}
}

func TestGetCustomerByName(t *testing.T) {
	var testCases = []struct {
		inp string
		out []customer
	}{

		{"Ishan Pandey", []customer{customer{1, "Ishan Pandey", "25-12-1999", address{1, "SSV Colony", "Lucknow", "UP", 1}}}},
		{"Varun Singh", []customer{{2, "Varun Singh", "23-12-1996", address{2, "St.Peters Street", "Singapore", "Malaysia", 2}}}},
		{"Shiva Tyagi", []customer(nil)},
	}
	for i := range testCases {
		req := httptest.NewRequest("GET", "/customer/"+testCases[i].inp, nil)
		w := httptest.NewRecorder()
		getCustomerByName(w, req)

		resp := w.Body.Bytes()
		var res []customer
		_ = json.Unmarshal(resp, &res)
		if w.Code != http.StatusOK {
			t.Logf("FAILED: ERROR IN CONNECTION")
		} else if !reflect.DeepEqual(res, testCases[i].out) {
			t.Errorf("FAILED for %v expected %v got %v\n", testCases[i].inp, testCases[i].out, res)
		}
	}
}

func TestGetCustomerById(t *testing.T) {
	var testCases = []struct {
		inp string
		out customer
	}{
		{"1", customer{1, "Ishan Pandey", "25-12-1999", address{1, "SSV Colony", "Lucknow", "UP", 1}}},
		{"2", customer{2, "Varun Singh", "23-12-1996", address{2, "St.Peters Street", "Singapore", "Malaysia", 2}}},
		// Empty struct.
		{"7", customer{}},
	}
	for i := range testCases {
		req := httptest.NewRequest("GET", "/customer/"+testCases[i].inp, nil)
		w := httptest.NewRecorder()
		getCustomerByID(w, req)

		resp := w.Body.Bytes()
		var res customer
		_ = json.Unmarshal(resp, &res)
		if w.Code != http.StatusOK {
			t.Logf("FAILED: ERROR IN CONNECTION")
		} else if !reflect.DeepEqual(res, testCases[i].out) {
			t.Errorf("FAILED for %v expected %v got %v\n", testCases[i].inp, testCases[i].out, res)
		}
	}
}

// In this we create customer. So, its necessary to be []byte.
func TestCreateCustomer(t *testing.T) {
	//var testCases = []struct {
	//	inp []byte
	//	// Ignore ID in the output, bcz we dont know the ID pre-handedly.
	//	out customer
	//}{
	//	//{[]byte("Name":"CustomerA","DOB":"10-10-2010","Addr.ID":"1","Addr.State":"Hyderabad","Addr.""Telangana", "12", 1),
	//	//	customer(3, "CustomerA", "10/10/2010", 1, "Hyderabad", "Telangana", "12", 1)},
	//	//
	//	//{"0",[]byte(``),[]Customer(nil)},
	//	//{"1234",[]byte(``),[]Customer(nil)},
	//	//
	//	//{customer{4, "Varun", "10-10-2011", address{2, "Patna", "Bihar", "121", 2}},
	//	//	customer{4, "CustomerB", "10/10/2011", address{2, "Patna", "Bihar", "121", 2}}},
	//	//
	//	//// ID 2 already exists.
	//	//{customer{2, "Akash", "10/10/2000", address{2, "Patna", "Bihar", "121", 2}},
	//	//	customer{}},
	//}

	//for i := range testCases {
	//	jsonInput, _ := json.Marshal(testCases[i].inp)
	//	req := httptest.NewRequest("POST", "/customer/", bytes.NewBuffer(jsonInput))
	//	w := httptest.NewRecorder()
	//	postCustomer(w, req)
	//
	//	resp := w.Body.Bytes()
	//	var res customer
	//	json.Unmarshal(resp, &res)
	//	if w.Code != http.StatusOK {
	//		t.Logf("FAILED: ERROR IN CONNECTION")
	//	} else if !reflect.DeepEqual(res, testCases[i].out) {
	//		t.Errorf("FAILED for %v Expected %v Got %v\n", testCases[i].inp, testCases[i].out, res)
	//	}
	//}

}

//}

func TestPutCustomer(t *testing.T) {
	var testCases = []struct {
		inp customer
		out customer
	}{
		// Changing name.
		{customer{ID: 1, Name: "Ritul Pandey", DOB: "", Addr: address{}},
			customer{1, "Ritul Pandey", "25-12-1999", address{1, "SSV Colony", "Lucknow", "UP", 1}}},
		// Trying to change DOB, not allowed
		{customer{ID: 2, Name: "", DOB: "10-10-1990", Addr: address{}},
			customer{2, "Varun Singh", "10-10-1990", address{2, "St.Peters Street", "Singapore", "Malaysia", 2}}},
		// Changing Address of ID that does not exist.
		{customer{ID: 100, Name: "CustomerB", DOB: "10/10/2011", Addr: address{2, "Patna", "Bihar", "121", 2}},
			customer{}},
		// Changing the address tha exists.
		{customer{2, "Varun Singh", "10-10-1990", address{2, "St.Peters Street", "Singapore", "SE Asia", 2}},
			customer{2, "Varun Singh", "10-10-1990", address{2, "St.Peters Street", "Singapore", "SE Asia", 2}}},
	}
	for i := range testCases {
		// Sending the input data in body of request.
		jsonInput, _ := json.Marshal(testCases[i].inp)
		req := httptest.NewRequest("POST", "/customer/", bytes.NewBuffer(jsonInput))
		w := httptest.NewRecorder()
		putCustomer(w, req)

		// Receiving JSON response data from response body.
		resp := w.Body.Bytes()
		var res customer
		_ = json.Unmarshal(resp, &res)
		if w.Code != http.StatusOK {
			t.Logf("FAILED: ERROR IN CONNECTION")
		} else if !reflect.DeepEqual(res, testCases[i].out) {
			t.Errorf("FAILED for %v Expected %v Got %v\n", testCases[i].inp, testCases[i].out, res)
		}
	}
}

func TestDeleteCustomer(t *testing.T) {
	var testCases = []struct {
		inpID string
		//out   customer
		out string
	}{
		{"1", "SUCCESS"},
		{"2", "SUCCESS"},

		{"100", ""},
		// Deletion output if it were customer object.
		//{"86", customer(nil)},
	}
	for i := range testCases {
		req := httptest.NewRequest("GET", "/customer/"+testCases[i].inpID, nil)
		w := httptest.NewRecorder()
		deleteCustomer(w, req)

		resp := w.Body.Bytes()
		var res string
		_ = json.Unmarshal(resp, &res)
		if w.Code != http.StatusOK {
			t.Logf("FAILED: ERROR IN CONNECTION")
		} else if !reflect.DeepEqual(res, testCases[i].out) {
			t.Errorf("FAILED for %v Expected %v Got %v\n", testCases[i].inpID, testCases[i].out, res)
		}
	}
}

/*package main

import (
	"bytes"
	"reflect"

	"io/ioutil"
	"net/http/httptest"
	"testing"
)

// RETRIEVE
// func TestGetCustomer tests func GetCustom(), which accepts request with /Customer/{name},
// to return all data (ID, NAME, ADDRESS) of all data in db.
// Query Parameter {name} is optional.
// Return an array of all customers.
func TestGetCustomer(t *testing.T) {
	testcases := []struct {
		input  string
		output [][]byte
	}{
		{"", [][]byte{
			[]byte(`{"ID" : 1, Name":"Ishan","Age":14,"Address":"India"}`),
			[]byte(`{"ID": 2, Name":"Suraj", "Age":18, "Address": "India"`),
		}},
		{"", [][]byte{
			[]byte(`{"ID" : 1, Name":"Ishan","Age":14,"Address":"India"}`),
			[]byte(`{"ID": 2, Name":"Suraj", "Age":18, "Address": "India"`),
		}},
	}

	t.Logf("Testing GetCustomer Func Handler: for no name or with name")
	for idx := range testcases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://localhost:8080/customer/"+testcases[idx].input, nil)
		GetCustomer(w, req)
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		// https://yourbasic.org/golang/compare-slices/

		// Convert a 2d slice to string:
		if reflect.DeepEqual(body, testcases[idx].output) {
			t.Error("Failed")
			//t.Logf("Expected: %s, \nGot %s", string(testcases[idx].output), string(body))
		}
	}
}

// RETRIEVE
// func TestGetCustomerByID tests func GetCustomerByID(), which accepts request with /Customer/{ID},
// to return all data belonging to a customer.
// Basically, returns a customer object.
func TestGetCustomerById(t *testing.T) {
	testcases := []struct {
		input  string
		output []byte
	}{
		{"", []byte(``)},
		{"", []byte(``)},
	}

	t.Logf("Testing GetCustomerByID Func Handler")
	for idx := range testcases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://localhost:8080/customer/"+testcases[idx].input, nil) //bytes.NewBuffer(testcases[idx].input))
		GetCustomerById(w, req)
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != string(testcases[idx].output) {
			t.Error("Failed")
			t.Logf("Expected: %s, \nGot %s", string(testcases[idx].output), string(body))
		}
	}

}

// POST: CREATE A RECORD IN DB (add a new record)
// func TestPostCustomer tests func POstCustomer(), which accepts request to add records of a customer
// given, {name, age, DOB} only if age >= 18.
// If the age is less than 18, then dont update any record.
func TestPostCustomer(t *testing.T) {
	testcases := []struct {
		input  []byte
		output []byte
	}{
		{[]byte(`{"Name":"Ishan","Age":14,"Address":"India"}`), []byte(`not eligible`)},
		{[]byte(`{"Name":"Ishan","Address":"India"}`), []byte(`not eligible`)},
		{[]byte(`{"Name":"XYZ","Age":21,"Address":"DNFKLEF"}`), []byte(`{"Name":"XYZ","Age":21,"Address":"DNFKLEF"}`)},
	}

	t.Logf("Testing PostCustomer Func Handler")
	for idx := range testcases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "http://localhost:8080/", bytes.NewBuffer(testcases[idx].input))
		PostCustomer(w, req)
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != string(testcases[idx].output) {
			t.Error("Failed")
			t.Logf("Expected: %s, \nGot %s", string(testcases[idx].output), string(body))
		}
	}
}

// PUT: EDIT/ UPDATE in DB.
// func TestPutCustomer tests PutCustomer(), which accepts /customer/{id, name/address}
// If DOB is passed, show error as we cant edit.
func TestPutCustomer(t *testing.T) {

}

// func TestDeleteCustomer tests DeleteCustomer(), which deletes /customer/{id}. in the database belonging to an ID.
// In the function, just return string "SUCCESS" in the implementation.
// We dont want to delete the data.
func TestDeleteCustomer(t *testing.T) {
	testcases := []struct {
		input  string
		output []byte
	}{
		{"1", []byte(`not eligible`)},
		{"2", []byte(`not eligible`)},
		{"3", []byte(`{"Name":"XYZ","Age":21,"Address":"DNFKLEF"}`)},
	}

	t.Logf("Testing PostCustomer Func Handler")
	for idx := range testcases {

	}
}

*/
