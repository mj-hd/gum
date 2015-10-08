package controllers

import "net/http"

func apiInit() {
}
func apiDel() {
}

func apiHandler(document http.ResponseWriter, request *http.Request) (status int, err error) {

	writeStruct(document, apiMember{
		Status:  "success",
		Message: "エラーなし",
	}, http.StatusOK)

	return http.StatusOK, nil
}
