package main

import (
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Middelware struct {
	writer  http.ResponseWriter
	request *http.Request
}

func (m Middelware) EnsureMethodIsAllowed(method string) {
	if m.request.Method != method {
		m.writer.WriteHeader(http.StatusMethodNotAllowed)
		m.writer.Write(m.MarshalToJson(model.ErrorResponse{
			Error: "Method is not allowed!",
		}))

		return
	}
}

func (m Middelware) MarshalToJson(value any) []byte {
	result, err := json.Marshal(value)

	if err != nil {
		http.Error(m.writer, err.Error(), http.StatusInternalServerError)
		return []byte{}
	}

	return result
}

func (m Middelware) SetResponseAsJson() {
	m.writer.Header().Set("Content-Type", "application/json")
}

func ArtifactBuilder() []string {
	artifact := []string{}
	lists, _ := ioutil.ReadFile("./data/list-study.txt")

	for _, value := range strings.Split(string(lists), "\n") {
		artifact = append(artifact, strings.ReplaceAll(value, "\r", ""))
	}

	return artifact
}

func UserListBuilder() []string {
	artifact := []string{}
	lists, _ := ioutil.ReadFile("./data/users.txt")

	for _, value := range strings.Split(string(lists), "\n") {
		artifact = append(artifact, strings.ReplaceAll(value, "\r", ""))
	}

	return artifact
}

func InOneDimentionOfArray(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func GetStudyProgram() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := Middelware{
			writer:  w,
			request: r,
		}

		m.SetResponseAsJson()
		m.EnsureMethodIsAllowed("GET")

		artifact := ArtifactBuilder()
		studyData := []model.StudyData{}

		for i := 0; i < len(artifact); i++ {
			builder := strings.Split(artifact[i], "_")

			studyData = append(studyData, model.StudyData{
				Code: builder[0],
				Name: builder[1],
			})
		}
		w.WriteHeader(http.StatusOK)
		w.Write(m.MarshalToJson(studyData))
	}
}

func AddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := Middelware{
			writer:  w,
			request: r,
		}

		m.SetResponseAsJson()
		m.EnsureMethodIsAllowed("POST")

		var user model.User

		userErr := json.NewDecoder(r.Body).Decode(&user)

		if userErr != nil {
			http.Error(w, userErr.Error(), http.StatusInternalServerError)
			return
		}
		if user.ID == "" || user.Name == "" || user.StudyCode == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(m.MarshalToJson(model.ErrorResponse{
				Error: "ID, name, or study code is empty",
			}))

			return
		}

		artifact := ArtifactBuilder()
		checker := []string{}

		for i := 0; i < len(artifact); i++ {
			checker = append(checker, strings.Split(artifact[i], "_")[0])
		}

		if !InOneDimentionOfArray(checker, user.StudyCode) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(m.MarshalToJson(model.ErrorResponse{
				Error: "study code not found",
			}))

			return
		}
		input := string(user.ID + "_" + user.Name + "_" + user.StudyCode + "\n")
		err := ioutil.WriteFile("data/users.txt", []byte(input), 0644)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(m.MarshalToJson(model.SuccessResponse{
			Username: user.ID,
			Message:  "add user success",
		}))
	}
}

func DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := Middelware{
			writer:  w,
			request: r,
		}

		m.SetResponseAsJson()
		m.EnsureMethodIsAllowed("DELETE")

		id := r.URL.Query().Get("id")

		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(m.MarshalToJson(model.ErrorResponse{
				Error: "user id is empty",
			}))

			return
		}

		userList := UserListBuilder()
		checker := []string{}

		for i := 0; i < len(userList); i++ {
			checker = append(checker, strings.Split(userList[i], "_")[0])
		}

		if !InOneDimentionOfArray(checker, id) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(m.MarshalToJson(model.ErrorResponse{
				Error: "user id not found",
			}))

			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(m.MarshalToJson(model.SuccessResponse{
			Username: id,
			Message:  "delete success",
		}))
	}
}

func main() {
	http.HandleFunc("/study-program", GetStudyProgram())
	http.HandleFunc("/user/add", AddUser())
	http.HandleFunc("/user/delete", DeleteUser())

	fmt.Println("starting web server at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
