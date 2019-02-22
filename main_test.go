package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// test health, nohealth

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize()

	code := m.Run()

	os.Exit(code)

}

// test the health URL returns something proper
func TestHealth(t *testing.T) {
	health_req, _ := http.NewRequest("GET", "/health", nil)
	health_rsp := executeRequest(health_req)

	checkResponseCode(t, http.StatusOK, health_rsp.Code)

	if body := health_rsp.Body.String(); body != `"ok"` {
		t.Errorf("Expected an 'ok' message. Got %s", body)
	} else {
		fmt.Println("- Test OK: health check")
	}
}

// test adding a new key
func TestAddKey(t *testing.T) {
	// Delete the keys folder, does it recreate it?
	os.RemoveAll(keysPath)
	// does it upload a key from a file?
	// does the key have a public key, does it alert if not?
	addkey_body, writer, addkey_upload_err := MultipartUpload("./test.pub")
  if addkey_upload_err != nil {
		t.Errorf("There was an error with the uploaded file.\nThis is NOT an application error but make sure the key exists")
  }
	addkey_req, _ := http.NewRequest("POST", "/addkey/?user=test", addkey_body)
	addkey_req.Header.Set("Content-Type", writer.FormDataContentType())
	addkey_rsp := executeRequest(addkey_req)

	checkResponseCode(t, http.StatusOK, addkey_rsp.Code)

	if body := addkey_rsp.Body.String(); body != `"Public key for user test saved"` {
		t.Errorf("Expected an empty message. Got %s", body)
	} else {
		fmt.Println("- Test OK: upload key without keys folder")
	}

	// does it save it to a file?
  if _, err := os.Stat("./keys/test.pub"); os.IsNotExist(err) {
		t.Errorf("Expected a new file. No file was written")
	} else {
		fmt.Println("- Test OK: writing key to local file")
  }

	// does the key have a username, does it alert if not?
	addkeynouser_body, writer, addkeynouser_upload_err := MultipartUpload("./test.pub")
  if addkeynouser_upload_err != nil {
		t.Errorf("There was an error with the uploaded file.\nThis is NOT an application error but make sure the key exists")
  }
	addkeynouser_req, _ := http.NewRequest("POST", "/addkey/", addkeynouser_body)
	addkeynouser_req.Header.Set("Content-Type", writer.FormDataContentType())
	addkeynouser_rsp := executeRequest(addkeynouser_req)

	checkResponseCode(t, http.StatusBadRequest, addkeynouser_rsp.Code)

	addkeynouser_expected := `"Error: no user was provided. Try ?user=username"`
	if body := addkeynouser_rsp.Body.String(); body != addkeynouser_expected {
		t.Errorf("Expected: "+addkeynouser_expected+". Got %s", body)
	} else {
		fmt.Println("- Test OK: upload key without username")
	}
	// given a testing pair, can it decrypt an encrypted message?
}

// test getting a key
func TestGetKey(t *testing.T) {
	// do we get a key?
	// given a testing pair, can it decrypt an encrypted message?
}

// test key registry

func TestKeyStore(t *testing.T) {
	// list key files vs keys on memory
	// remove a key
	// clean up all keys
}

//https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func MultipartUpload(path string) (*bytes.Buffer, *multipart.Writer, error) {
	paramName := "file"
	params := map[string]string{}
  var out_err error

	file, err := os.Open(path)
	if err != nil {
    out_err = err
		//panic(err)
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		fmt.Println(err)
	}
	_, err = io.Copy(part, file)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
	}
	return body, writer, out_err
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
