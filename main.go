package main

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
)

type task struct {
	Id int `json:"id"`
	TaskName string `json:"task_name,omitempty"`
	TaskContent string `json:"task_content,omitempty"`
}

type result struct {
	Code int `json:"code"`
	Data []task `json:"data"`
	Message string `json:"message"`
}

var primaryKey = 0

var taskList []task
var taskListMap = make(map[int]task)


func httpInterceptor(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")

		fn(w, r)
	}
}

/**
Create a new task
 */
func onAddTask(w http.ResponseWriter, r *http.Request) {
	var newTask map[string]string
	res, _ := json.Marshal(result{Code: 200, Message: "新增成功"})

	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		w.WriteHeader(http.StatusBadGateway)
		errorRes, _ := json.Marshal(result{Code: 500, Message: "服务器内部错误"})
		w.Write(errorRes)
	}

	taskList = append(taskList, task{Id: primaryKey, TaskName: newTask["task_name"], TaskContent: newTask["task_content"]})
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

/**
Delete a task
 */
func onDeleteTask(w http.ResponseWriter, r *http.Request) {
	var taskInfo = task{}

	if err := json.NewDecoder(r.Body).Decode(&taskInfo); err != nil {
		errRes, _ := json.Marshal(result{
			Code: http.StatusBadRequest,
			Message: "缺少参数对象",
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(errRes)
		return
	}
	
	if !reflect.ValueOf(taskInfo).FieldByName("Id").IsValid() {
		errRes, _ := json.Marshal(result{
			Code:    http.StatusBadRequest,
			Message: "缺少必填参数 id",
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(errRes)
		return
	}

	taskId := taskInfo.Id

	 if _, ok := taskListMap[taskId]; !ok  {
		errRes, _ := json.Marshal(result{
			Code: http.StatusBadRequest,
			Message: "没有找到相对应的任务",
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(errRes)
		return
	}

	delete(taskListMap, taskId)
	for index, task := range taskList {
		if task.Id == taskId {
			taskList = append(taskList[:index], taskList[index+1:]...)
			break
		}
	}

	res, _ := json.Marshal(result{
		Code: 200,
		Message: "删除成功",
	})

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

/**
Update a existed task
 */
func onUpdateTask(w http.ResponseWriter, r *http.Request) {
	var wannaUpdateTask = task{}
	var taskId = wannaUpdateTask.Id

	if err := json.NewDecoder(r.Body).Decode(&wannaUpdateTask); err != nil {
		errRes, _ := json.Marshal(result{
			Code: http.StatusInternalServerError,
			Message: "服务器内部错误",
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errRes)
		return
	}

	taskList[taskId] = wannaUpdateTask
	taskListMap[taskId] = wannaUpdateTask

	res, _ := json.Marshal(result{
		Code: http.StatusBadRequest,
		Message: "更新成功",
	})

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// Get task list
func getTaskList(w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(result{Code: 200, Data: taskList})

	if err != nil {
		log.Fatal(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// Default task list
func initTaskList() {
	for i := 0; i < 3; i++ {
		var task = task{Id: primaryKey, TaskName: "任务", TaskContent: "今天要努力的工作"}
		taskList = append(taskList, task)
		taskListMap[primaryKey] = task
		primaryKey = primaryKey + 1
	}
}


func main() {
	initTaskList()

	http.HandleFunc("/task/create", httpInterceptor(onAddTask))
	http.HandleFunc("/task/delete", httpInterceptor(onDeleteTask))
	http.HandleFunc("/task/update", httpInterceptor(onUpdateTask))
	http.HandleFunc("/task/get", httpInterceptor(getTaskList))

	http.ListenAndServe(":50000", nil)
}

