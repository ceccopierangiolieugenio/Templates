Code Test
=================

## Summary
We need an HTTP server that will accept any POST request (JSON) from multiple clients' websites. Each request forms part of a struct (for that particular visitor) that will be printed to the terminal when the struct is fully complete. 

## Frontend (JavaScript)
Insert JavaScript into `index.html` (supplied) that captures and posts data every time one of the below events happens; this means you will be posting multiple times per visitor. 

  - if the screen resizes, the before and after dimensions, assume only one resize occurs
  - copy & paste (for each field)
  - time taken, in seconds, from the first character being typed to submitting the form

### Example JSON Requests
```javascript
{
  "eventType": "copyAndPaste",
  "websiteUrl": "https://github.com/ceccopierangiolieugenio",
  "sessionId": "123123-123123-123123123",
  "pasted": true,
  "formId": "inputCardNumber"
}

{
  "eventType": "timeTaken",
  "websiteUrl": "https://github.com/ceccopierangiolieugenio",
  "sessionId": "123123-123123-123123123",
  "timeTaken": 72,
}

```

## Backend (Go)

The backend should:

1. Create a server
2. Accept POST requests in JSON format similar to those specified above
3. Map the JSON requests to relevant sections of a Event struct (specified below)
4. Print the struct for each stage of its construction
5. Also print the struct when it is complete (i.e. when the form has been submitted)

We would like the server to be written to handle multiple requests arriving on
the same session at the same time. We'd also like to see some tests.


### Go Struct
```go
type Event struct {
	WebsiteUrl         string
	SessionId          string
	ResizeFrom         Dimension
	ResizeTo           Dimension
	CopyAndPaste       map[string]bool // map[fieldId]true
	FormCompletionTime int // Seconds
}

type Dimension struct {
	Width  string
	Height string
}
```
