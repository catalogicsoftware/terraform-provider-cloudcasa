package handler

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
)

type CreateKubeclusterReq struct {
	Name				string		`json:"name"`
	//Cc_user_email		string		`json:"cc_user_email"`
	//Org_id				string		`json:"org_id"`
	// Status 				struct {}	`json:"status"`
	// Links 				struct {}	`json:"_links"`
}

type CreateKubeclusterResp struct {
	Id					string		`json:"_id"`
	Name				string		`json:"name"`
	Cc_user_email		string		`json:"cc_user_email"`
	Updated				string		`json:"_updated"`
	Created				string		`json:"_created"`
	Etag				string		`json:"_etag"`
	Org_id				string		`json:"org_id"`
	Status 				string		`json:"_status"`
	Links 				struct {}	`json:"_links"`
}

type GetKubeclusterReq struct {
	Id					string		`json:"_id"`
}

type KubeclusterStatus struct {
	State				string		`json:"state"`
	Agent_url		    string		`json:"agentUrl"`
}

type GetKubeclusterResp struct {
	Id					string				`json:"_id"`
	Name				string				`json:"name"`
	Cc_user_email		string				`json:"cc_user_email"`
	Updated				string				`json:"_updated"`
	Created				string				`json:"_created"`
	Etag				string				`json:"_etag"`
	Org_id				string				`json:"org_id"`
	Status 				KubeclusterStatus	`json:"status"`
	Links 				struct {}			`json:"_links"`
} 

// Create kubecluster resource using cloudcasa API
func CreateKubecluster(CreateKubeclusterReq CreateKubeclusterReq) *CreateKubeclusterResp {
	// Create rest request struct
	kubeclusterCreate, _ := json.Marshal(CreateKubeclusterReq)

	url := "https://api.staging.cloudcasa.io/api/v1/kubeclusters/"
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InM3NmtuNThRT2liTXRfZnNpVFlLMCJ9.eyJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9jb3VudHJ5IjoiVW5pdGVkIFN0YXRlcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL3RpbWV6b25lIjoiQW1lcmljYS9OZXdfWW9yayIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvdW50cnlfY29kZSI6IlVTIiwiaHR0cDovL3d3dy5jbG91ZGNhc2EuaW8vY291bnRyeV9jb2RlMyI6IlVTQSIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2ZpcnN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9sYXN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9qb2JUaXRsZSI6IkRldm9wcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvbXBhbnkiOiJDYXRhbG9naWMgU29mdHdhcmUiLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9hd3NfbWFya2V0cGxhY2VfdG9rZW4iOiItIiwibmlja25hbWUiOiJqZ2FybmVyIiwibmFtZSI6IkpvbmF0aGFuIEdhcm5lciIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci8yOTlhNmJhNjhlNjEwOGFiYjY1MmY4ZTkwZTM0YjVhNj9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmpnLnBuZyIsInVwZGF0ZWRfYXQiOiIyMDIzLTAyLTI0VDE4OjU0OjEzLjMxN1oiLCJlbWFpbCI6ImpnYXJuZXJAY2F0YWxvZ2ljc29mdHdhcmUuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImlzcyI6Imh0dHBzOi8vYXV0aC5jbG91ZGNhc2EuaW8vIiwiYXVkIjoiSkVKU3plblBGeE5FUFEwaDY0ZDIzZTZRMEdKNXpRanQiLCJpYXQiOjE2NzcyNjQ4NTQsImV4cCI6MTY3NzI3MjA1NCwic3ViIjoiYXV0aDB8NWZhYzQ4NDg0MWQ3MDgwMDY4YTA2ZGM5Iiwic2lkIjoid1E4NDg2TUhzd2dERjIzdVh0OGhyX1VrZ2hKT2x6cDIiLCJub25jZSI6ImJFRXlhRE5JZDJWUmJtdHVORVpOZDNWTlNtSllZMEZTYlVweFpIWkRTRk5OZVhCS2RsUTNkakZhZEE9PSJ9.iwYbSSSlR3l8oNx97NQVdymfWnc4iuhgmuG3XNv19vQqx5OKyX5KCZBdf_KUR4BxGtB6YfBV1BAJrsuG_8UGrRYghvBe6mBFKetFZ-Wm0_eu8O1dn1mnp7HsX-llZOMU78jarfpZWoLYCRh8M8WXYOSlFJlxdFW1HyhFuhGVX426uQCIuWUfv2M99XHzM2vKmdU-RI8FRQFsJgkrrETRlIrXO9ANBk-f5nkZLoGiigkUeINXHk75kf4AxJVpeoXl8dtqLbkq3G0WjYaeX8T95IK9wx5CzLvfrBuQ0VK5FNgoDATIuJpwSDTO_rt0GPjrIz8Pt-96a3AtBadtm30KMA"

	// makeHttpRequest is in handler/client.go
	respBody := makeHttpRequest(url, http.MethodPost, JSON, kubeclusterCreate, token)

	var apiCreateKubeclusterResp CreateKubeclusterResp
	json.Unmarshal(respBody, &apiCreateKubeclusterResp)
	return &apiCreateKubeclusterResp
}

func GetKubecluster(kubeclusterId string) *GetKubeclusterResp {
	// Create rest request struct
	//kubeclusterGet, _ := json.Marshal(GetKubeclusterReq)

	url := "https://api.staging.cloudcasa.io/api/v1/kubeclusters/"
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InM3NmtuNThRT2liTXRfZnNpVFlLMCJ9.eyJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9jb3VudHJ5IjoiVW5pdGVkIFN0YXRlcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL3RpbWV6b25lIjoiQW1lcmljYS9OZXdfWW9yayIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvdW50cnlfY29kZSI6IlVTIiwiaHR0cDovL3d3dy5jbG91ZGNhc2EuaW8vY291bnRyeV9jb2RlMyI6IlVTQSIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2ZpcnN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9sYXN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9qb2JUaXRsZSI6IkRldm9wcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvbXBhbnkiOiJDYXRhbG9naWMgU29mdHdhcmUiLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9hd3NfbWFya2V0cGxhY2VfdG9rZW4iOiItIiwibmlja25hbWUiOiJqZ2FybmVyIiwibmFtZSI6IkpvbmF0aGFuIEdhcm5lciIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci8yOTlhNmJhNjhlNjEwOGFiYjY1MmY4ZTkwZTM0YjVhNj9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmpnLnBuZyIsInVwZGF0ZWRfYXQiOiIyMDIzLTAyLTI0VDE4OjU0OjEzLjMxN1oiLCJlbWFpbCI6ImpnYXJuZXJAY2F0YWxvZ2ljc29mdHdhcmUuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImlzcyI6Imh0dHBzOi8vYXV0aC5jbG91ZGNhc2EuaW8vIiwiYXVkIjoiSkVKU3plblBGeE5FUFEwaDY0ZDIzZTZRMEdKNXpRanQiLCJpYXQiOjE2NzcyNjQ4NTQsImV4cCI6MTY3NzI3MjA1NCwic3ViIjoiYXV0aDB8NWZhYzQ4NDg0MWQ3MDgwMDY4YTA2ZGM5Iiwic2lkIjoid1E4NDg2TUhzd2dERjIzdVh0OGhyX1VrZ2hKT2x6cDIiLCJub25jZSI6ImJFRXlhRE5JZDJWUmJtdHVORVpOZDNWTlNtSllZMEZTYlVweFpIWkRTRk5OZVhCS2RsUTNkakZhZEE9PSJ9.iwYbSSSlR3l8oNx97NQVdymfWnc4iuhgmuG3XNv19vQqx5OKyX5KCZBdf_KUR4BxGtB6YfBV1BAJrsuG_8UGrRYghvBe6mBFKetFZ-Wm0_eu8O1dn1mnp7HsX-llZOMU78jarfpZWoLYCRh8M8WXYOSlFJlxdFW1HyhFuhGVX426uQCIuWUfv2M99XHzM2vKmdU-RI8FRQFsJgkrrETRlIrXO9ANBk-f5nkZLoGiigkUeINXHk75kf4AxJVpeoXl8dtqLbkq3G0WjYaeX8T95IK9wx5CzLvfrBuQ0VK5FNgoDATIuJpwSDTO_rt0GPjrIz8Pt-96a3AtBadtm30KMA"

	respBody := makeHttpRequest(url + kubeclusterId, http.MethodGet, JSON, nil, token)
	
	var apiGetKubeclusterResp GetKubeclusterResp
	json.Unmarshal(respBody, &apiGetKubeclusterResp)

	file, _ := json.MarshalIndent(respBody, "", " ")
	_ = ioutil.WriteFile("get_resp.json", file, 0644)

	return &apiGetKubeclusterResp
}

