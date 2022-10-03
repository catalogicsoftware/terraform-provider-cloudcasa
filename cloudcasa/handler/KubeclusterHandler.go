package handler

import (
	"encoding/json"
	"net/http"
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

func CreateKubecluster(CreateKubeclusterReq CreateKubeclusterReq) *CreateKubeclusterResp {
	kubeclusterCreate, _ := json.Marshal(CreateKubeclusterReq)
	url := "https://home.cloudcasa.io/api/v1/kubeclusters"
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InM3NmtuNThRT2liTXRfZnNpVFlLMCJ9.eyJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9jb3VudHJ5IjoiVW5pdGVkIFN0YXRlcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL3RpbWV6b25lIjoiQW1lcmljYS9OZXdfWW9yayIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvdW50cnlfY29kZSI6IlVTIiwiaHR0cDovL3d3dy5jbG91ZGNhc2EuaW8vY291bnRyeV9jb2RlMyI6IlVTQSIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2ZpcnN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9sYXN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9qb2JUaXRsZSI6IkRldm9wcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvbXBhbnkiOiJDYXRhbG9naWMgU29mdHdhcmUiLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9hd3NfbWFya2V0cGxhY2VfdG9rZW4iOiItIiwibmlja25hbWUiOiJqZ2FybmVyIiwibmFtZSI6IkpvbmF0aGFuIEdhcm5lciIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci8yOTlhNmJhNjhlNjEwOGFiYjY1MmY4ZTkwZTM0YjVhNj9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmpnLnBuZyIsInVwZGF0ZWRfYXQiOiIyMDIyLTEwLTAzVDE3OjM3OjAxLjg2N1oiLCJlbWFpbCI6ImpnYXJuZXJAY2F0YWxvZ2ljc29mdHdhcmUuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImlzcyI6Imh0dHBzOi8vYXV0aC5jbG91ZGNhc2EuaW8vIiwic3ViIjoiYXV0aDB8NWZhYzQ4NDg0MWQ3MDgwMDY4YTA2ZGM5IiwiYXVkIjoiSkVKU3plblBGeE5FUFEwaDY0ZDIzZTZRMEdKNXpRanQiLCJpYXQiOjE2NjQ4MTg2MjQsImV4cCI6MTY2NDgyNTgyNCwic2lkIjoiYU5xOXdNRXdRaS03by1qa2drZktac0l1WkVZZjJ5dnoiLCJub25jZSI6IkxtOHRUa3hvVmtsamFsSTNkRlJCVXpWM1JWaG9Mak14WlVjMlNEaElXRXh3VDNSVVMxQlZZVFZ5Y1E9PSJ9.TXsrznqtqw0Jx2YIllTaCdIz--Kf3fqpkDpkj_FqicaUX-ZHQwU57Qu95subO8thzhku9o5Nw_twWCH3JU7HAGb2MH9biVOLvHjf3bkGbHMLipAS8REjPFRRsroWE5AzcceGKQf_VttciPTR8_mZrbTRMdv7oMpLyE92m8eVCP8_41avKCj0UdvtW8FSjNXyTBwl2NCTYz4Ubo4H4mt5LOYBulIoKr8A1mx2IDTjYRBGapLwisno9hS_dGjUD7T0deXe0LZ0od8vv92-7IUyfYXiAuPYpvkIO-Nv8YsPrnkRDxD6vFoJ19CCj0PE9BkbaxiMYIk2tOCK3NYFRCytsw"

	// makeHttpRequest is in handler/client.go
	respBody := makeHttpRequest(url, http.MethodPost, JSON, kubeclusterCreate, token)

	var apiCreateKubeclusterResp CreateKubeclusterResp
	json.Unmarshal(respBody, &apiCreateKubeclusterResp)
	return &apiCreateKubeclusterResp
}



