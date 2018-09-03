package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// helper to set up testing json file and tell the app to use it
func setupDataFileForTest() {
	copy("data.json.example", "test.json")
	setupGlobalOptions("test.json")
}

// this is ridiculous. this is how to copy a file in go.
func copy(src string, dest string) error {
	in, err := os.Open(src)

	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func TestRootRoute(t *testing.T) {
	setupDataFileForTest()
	router := setupAPI()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Response was incorrect, got: %d, want: 200.", w.Code)
	}
}

func TestCollectionRoute(t *testing.T) {
	setupDataFileForTest()
	router := setupAPI()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/notes", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Response was incorrect, got: %d, want: 200.", w.Code)
	}
}

func TestCollectionRouteWhenNoneFound(t *testing.T) {
	setupDataFileForTest()
	router := setupAPI()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/derp", nil)
	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Response was incorrect, got: %d, want: 404.", w.Code)
	}
}

func TestCollectionRecordRoute(t *testing.T) {
	setupDataFileForTest()
	router := setupAPI()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/notes/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Response was incorrect, got: %d, want: 200.", w.Code)
	}
}

func TestCollectionRecordRouteForNonExistantID(t *testing.T) {
	setupDataFileForTest()
	router := setupAPI()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/notes/12345", nil)
	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Response was incorrect, got: %d, want: 404.", w.Code)
	}
}

func TestCreateRecord(t *testing.T) {
	setupDataFileForTest()
	router := setupAPI()

	var jsonStr = []byte(`{"title":"new record"}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/notes", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	if w.Code != 201 {
		t.Errorf("Response was incorrect, got: %d, want: 201.", w.Code)
	}
}

func TestUpdateRecordViaPut(t *testing.T) {
	setupDataFileForTest()
	router := setupAPI()

	var jsonStr = []byte(`{"title":"new record"}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/notes/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Response was incorrect, got: %d, want: 200.", w.Code)
	}
}

func TestUpdateRecordViaPatch(t *testing.T) {
	setupDataFileForTest()
	router := setupAPI()

	var jsonStr = []byte(`{"title":"new record"}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/notes/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Response was incorrect, got: %d, want: 200.", w.Code)
	}
}

// unit tests

func TestFindByIdWithExistingRecord(t *testing.T) {

	setupDataFileForTest()
	_, err := findById("notes", "1")

	if err != nil {
		t.Errorf("Record exists but was not found")
	}
}

func TestFindByIdWithNonInt(t *testing.T) {

	setupDataFileForTest()
	_, err := findById("notes", "ABCD")

	if err == nil {
		t.Errorf("Not getting an error when 'ABCD' used as key.")
	}
}

func TestFindByIdWithNonExistingRecord(t *testing.T) {
	setupDataFileForTest()

	_, err := findById("notes", "1234")

	if err == nil {
		t.Errorf("Should be error when doesn't exist")
	}
}

func TestRemoveRecordByIdWorks(t *testing.T) {

	setupDataFileForTest()

	removeById("notes", "1")

	// see if it persisted
	results, _ := getData()
	children, _ := results.S("notes").Children()

	for _, child := range children {

		if child.S("id").Data().(float64) == 1 {
			t.Errorf("New record still found in the json - failed")
			break
		}
	}

}

func TestRemoveRecordByIdWithInvalidID(t *testing.T) {

	setupDataFileForTest()

	_, err := removeById("notes", "1234")

	if err == nil {
		t.Errorf("Should be error when doesn't exist")
	}

}

func TestCreateNewRecord(t *testing.T) {
	setupDataFileForTest()

	var result interface{}
	createNewRecord("notes", []byte(`{"title" : "test"}`))

	// see if it persisted
	results, _ := getData()
	children, _ := results.S("notes").Children()

	// find the index of the record we have to delete
	for _, child := range children {

		// if we find it....
		if child.S("title").Data().(string) == "test" {
			// save the record we found as the result along with the index
			result = child
			break
		}
	}

	if result == nil {
		t.Errorf("New record wasn't found in the json - failed")
	}
}

func TestUpdateRecordAtID(t *testing.T) {

	setupDataFileForTest()

	var result = false

	updateById("notes", "1", []byte(`{"title" : "blah blah blah"}`))

	// see if it persisted
	records, _ := getData()
	children, _ := records.S("notes").Children()

	// find the index of the record we have to delete
	for _, child := range children {

		// if we find it....
		if child.S("id").Data().(float64) == 1 && child.S("title").Data().(string) == "blah blah blah" {
			// save the record we found as the result along with the index
			result = true
			break
		}
	}

	if !result {
		t.Errorf("New record wasn't found in the json - failed")
	}
}

func TestUpdateRecordAtBadID(t *testing.T) {

	setupDataFileForTest()

	_, err := updateById("notes", "1234", []byte(`{"title" : "blah blah blah"}`))

	if err == nil {
		t.Errorf("Update for ID 1234 should have failed but didn't")
	}
}

func TestUpdateRecordAtinvalidID(t *testing.T) {

	setupDataFileForTest()

	_, err := updateById("notes", "ABCD", []byte(`{"title" : "blah blah blah"}`))

	if err == nil {
		t.Errorf("Update for ID ABCD should have failed but didn't")
	}
}
