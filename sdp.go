package sdp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type SDP struct {
	Server string
	Port   string
	APIKey string
}

type SDPResponse struct {
	Operation Operation `json:"operation"`
}

type Operation struct {
	Result  Result  `json:"result"`
	Details Details `json:"Details"`
}

type Details struct {
	WORKORDERID  string
	CustomFields map[string]interface{}
}

type Result struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

//Reply area
type SDPReplyResponse struct {
	Operation ReplyOperation `json:"operation"`
}

type ReplyOperation struct {
	Result ReplyResult `json:"result"`
}

type ReplyResult struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

///////////

//NewTask - create new task to ServiceDesk plus
func (s SDP) NewTask(subject, body string) (string, error) {
	//http://10.199.1.174:8080/sdpapi/request
	/*
	   "url": "http://10.199.1.174:8080/sdpapi/request",
	    "headers": {
	      "Accept": "application/json",
	      "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
	      "TECHNICIAN_KEY": "539F22CF-5F41-4928-A726-D23E1695659F"
	    },
	    "data": "format=json&OPERATION_NAME=ADD_REQUEST&INPUT_DATA={\n    \"operation\": {\n        \"details\": {\n            \"requester\": \"Vasiliy Terkin\",\n            \"subject\": \"API REQUEST  #009\",\n            \"description\": \"Specify Description and many orther words\"\n        }\n    }\n}",
	    "timeout": {}
	*/
	uri := "http://" + s.Server + ":" + s.Port + "/sdpapi/request"

	//Replace & from body request
	body = strings.Replace(body, "&", "+", -1) //Replace & - don'twork

	fmt.Println(`!!!!BODY REQUEST:format=json&OPERATION_NAME=ADD_REQUEST&INPUT_DATA={"operation": {"details": {"requester": "Watcom robot","subject": "` + subject + `","description": "` + body + `","priority": "High"}}}`)

	bodyrequest := []byte(`format=json&OPERATION_NAME=ADD_REQUEST&INPUT_DATA={"operation": {"details": {"requester": "Watcom robot","subject": "` + subject + `","description": "` + body + `","priority": "High"}}}`)

	//Set reader from body request
	r := bytes.NewReader(bodyrequest)
	//1 Check same task - search by key "field1097": ""
	client := &http.Client{}
	req, err := http.NewRequest(
		"POST", uri, r,
	)

	// добавляем заголовки
	req.Header.Add("Accept", "application/json")                                       // добавляем заголовок Accept
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8") // добавляем заголовок Content-Type
	req.Header.Add("TECHNICIAN_KEY", s.APIKey)                                         // добавляем заголовок Content-Type

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result SDPResponse //Response code>200, then error hendler
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Println(err.Error())
		}

	}

	//All is OK< web response code is 200 OK!
	var result SDPResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		//return res, err
		fmt.Println(err.Error())
	}
	fmt.Printf("result raw=%v \n", result)
	taskid := result.Operation.Details.WORKORDERID
	return taskid, nil
}

//
//NewTask - create new task to ServiceDesk plus
func (s SDP) ReplyTask(taskID, toAddress, subject, body string) (bool, error) {
	/*

		  "method": "POST",
		  "transformRequest": [
		    null
		  ],
		  "transformResponse": [
		    null
		  ],
		  "url": "http://10.199.1.174:8080/sdpapi/request/35",
		  "headers": {
		    "Accept": "application/json, text/plain, ",
		    "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
		    "TECHNICIAN_KEY": "539F22CF-5F41-4928-A726-D23E1695659F"
		  },
		  "data": "format=json&OPERATION_NAME=REPLY_REQUEST&INPUT_DATA={\n    \"operation\": {\n        \"details\": {\n
			\"to\": \"devnull@watcom.ru\",\n
			\"subject\": \"Add Comment -02\",\n
			\"description\": \"Body comment\"\n        }\n    }\n}",
		  "timeout": {}
		}
	*/
	uri := "http://" + s.Server + ":" + s.Port + "/sdpapi/request/" + taskID

	//Replace & from body request
	body = strings.Replace(body, "&", "+", -1) //Replace & - don'twork

	bodyrequest := []byte(`format=json&OPERATION_NAME=REPLY_REQUEST&INPUT_DATA={"operation": {"details": {"to": "` + toAddress + `","subject": "` + subject + `","description": "` + body + `"}}}`)

	//Set reader from body request
	r := bytes.NewReader(bodyrequest)
	//1 Check same task - search by key "field1097": ""
	client := &http.Client{}
	req, err := http.NewRequest(
		"POST", uri, r,
	)

	// добавляем заголовки
	req.Header.Add("Accept", "application/json")                                       // добавляем заголовок Accept
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8") // добавляем заголовок Content-Type
	req.Header.Add("TECHNICIAN_KEY", s.APIKey)                                         // добавляем заголовок Content-Type

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result SDPReplyResponse //Response code>200, then error hendler
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Println(err.Error())
			return false, err
		}
		err := errors.New("Http StatusCode: " + fmt.Sprintf("%d, result: %+v", resp.StatusCode, result))
		return false, err
	}

	//All is OK< web response code is 200 OK!
	var result SDPReplyResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return true, err
	}
	fmt.Printf("result raw=%+v \n", result)
	//taskid := result.Operation.Details.WORKORDERID
	return true, nil
}
