package main

import (
    "log"
    "flag"
    "testing"
    "bytes"
    "reflect"
    "io/ioutil"
	"net/http"
    "net/http/httptest"
    "encoding/json"
)

var clientdir string

func init() {
    flag.StringVar(&clientdir, "clientdir", "./client", "Client Directory")
}

func PerformGET(t *testing.T, uri string, file string) {
    log.Printf("[GET] Testing, uri:%s , file:%s", uri, file)
    req, err := http.NewRequest("GET", uri, nil)
    if err != nil {
        t.Fatal(err)
    }

    byt, err := ioutil.ReadFile(file)
    if err != nil {
        t.Fatal(err)
    }
    content := string(byt)

    srv := &Server{ Events:    make(map[string]*Event),
                    ClientDir: clientdir}

    // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(srv.CallbackHandler)

    // Our handlers satisfy http.Handler, so we can call their ServeHTTP method 
    // directly and pass in our Request and ResponseRecorder.
    handler.ServeHTTP(rr, req)

    // Check the response body is what we expect.
    if rr.Body.String() != content {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), content)
    }
}

func TestGET(t *testing.T) {
    PerformGET(t, "/", clientdir+"/index.html")
    PerformGET(t, "/js/client.js", clientdir+"/js/client.js")
}

func PerformPOST(t *testing.T, srv *Server, uri string, data interface{}, content Event, expectedStatus int) {
    body, _ := json.Marshal(data)
    //req, err := http.NewRequest("POST", "/questions/", body)
    //var m map[string]interface{}
    //err = json.NewDecoder(req.Body).Decode(&m)
    //req.Body.Close()
    //fmt.Println(err, m)


    req, err := http.NewRequest("POST", uri, bytes.NewReader(body))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Add("Content-Type", "application/json")

    // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(srv.CallbackHandler)

    // Our handlers satisfy http.Handler, so we can call their ServeHTTP method 
    // directly and pass in our Request and ResponseRecorder.
    handler.ServeHTTP(rr, req)

    // Check the status code is what we expect.
    if status := rr.Code; status != expectedStatus {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, expectedStatus)
    }

    if expectedStatus != http.StatusOK {
        /* If a testing error code is performed
           there is no need to compare the return content */
        return
    }

    if !reflect.DeepEqual(*srv.Events[content.SessionId], content) {
        t.Errorf("handler returned wrong content: got %v want %v",*srv.Events[content.SessionId], content)
    }
}

func TestPOSTValid(t *testing.T) {
    srv := &Server{ Events:    make(map[string]*Event),
        ClientDir: clientdir}

    var data   interface{}
    var expected Event

     /* Test a CopyAndPast callback */
    data = map[string]interface{}{
        "eventType": "copyAndPaste",
        "websiteUrl": "https://github.com/ceccopierangiolieugenio",
        "sessionId": "123123-123123-123123123",
        "pasted": true,
        "formId": "inputCardNumber",
      }
    expected = Event{
        WebsiteUrl: "https://github.com/ceccopierangiolieugenio",
        SessionId: "123123-123123-123123123",
        ResizeFrom: Dimension{Width: "", Height: ""},
        ResizeTo:   Dimension{Width: "", Height: ""},
        CopyAndPaste: map[string]bool{"inputCardNumber": true},
        FormCompletionTime:-1,
    }
    PerformPOST(t, srv, "/", data, expected, http.StatusOK)

    /* Test a second CopyAndPast callback for the same sessionId 
       the should update the previous struct */
    data = map[string]interface{}{
        "eventType": "copyAndPaste",
        "websiteUrl": "https://github.com/ceccopierangiolieugenio",
        "sessionId": "123123-123123-123123123",
        "pasted": false,
        "formId": "inputEmail",
      }
    expected.CopyAndPaste["inputEmail"] = false
    PerformPOST(t, srv, "/", data, expected, http.StatusOK)
}

func TestPOSTInvalid(t *testing.T) {
    srv := &Server{ Events:    make(map[string]*Event),
        ClientDir: clientdir}

    var data   interface{}
    var expected Event

    /* Test an invalid input */
    data = map[string]interface{}{}
    PerformPOST(t, srv, "/", data, expected, http.StatusBadRequest)

}


func TestPrettyPint(t *testing.T) {
    rr := httptest.NewRecorder()
    expected := Event{
        WebsiteUrl: "https://github.com/ceccopierangiolieugenio",
        SessionId: "123123-123123-123123123",
        FormCompletionTime:-1,
    }
    PrettyPrint(&expected, rr)
}
