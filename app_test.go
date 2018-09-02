package main

import (
	"fmt"
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
func copy(src string, dest string) {
	from, _ := os.Open(src)
	defer from.Close()

	to, _ := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0666)

	defer to.Close()

	io.Copy(to, from)
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

	result, err := findById("notes", "1234")

	fmt.Print(result.String())
	if err == nil {
		t.Errorf("Should be error when doesn't exist")
	}
}

func TestRemoveRecordByIdWorks(t *testing.T) {

	setupDataFileForTest()

	result, _ := removeById("notes", "1")

	fmt.Print(result.S("id").Data().(float64))

}

func TestRemoveRecordByIdWithInvalidID(t *testing.T) {

	setupDataFileForTest()

	_, err := removeById("notes", "1234")

	if err == nil {
		t.Errorf("Should be error when doesn't exist")
	}

}

func TestCreateNewRecord(t *testing.T) {

	var result interface{}
	setupDataFileForTest()
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
