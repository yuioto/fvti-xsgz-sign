package savestusignin

import "fmt"

func GetTaskId(name string) (string, error) {
	id := "writting"

	taskList := GetTaskList()

	fmt.Println(taskList)

	return id, nil
}

func GetTaskList() (string){
	return "writting"
}