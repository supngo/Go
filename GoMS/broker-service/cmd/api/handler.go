package main

import (
	"broker/event"
	"broker/logs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(res http.ResponseWriter, req *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(res, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(res http.ResponseWriter, req *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(res, req, &requestPayload)

	if err != nil {
		app.errorJSON(res, err)
		return

	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(res, requestPayload.Auth)
	case "log":
		// app.logItem(res, requestPayload.Log)
		// app.logEventViaRabbit(res, requestPayload.Log)
		app.logItemViaRPC(res, requestPayload.Log)
	default:
		app.errorJSON(res, errors.New("unknown action"))
	}
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) authenticate(res http.ResponseWriter, load AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(load, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(res, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(res, err)
		return
	}

	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(res, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(res, errors.New("error calling auth service"))
		return
	}

	// create a varabiel we'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)

	if err != nil {
		app.errorJSON(res, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(res, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data
	app.writeJSON(res, http.StatusAccepted, payload)
}

// logEventViaRabbit logs an event using the logger-service. It makes the call by pushing the data to RabbitMQ.
func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	app.writeJSON(w, http.StatusAccepted, payload)
}

// pushToQueue pushes a message into RabbitMQ
func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}

type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) logItemViaRPC(res http.ResponseWriter, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJSON(res, err)
		return
	}

	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJSON(res, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	app.writeJSON(res, http.StatusAccepted, payload)
}

func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(w, http.StatusAccepted, payload)
}
