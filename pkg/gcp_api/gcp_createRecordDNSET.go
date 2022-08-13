package gcp_api

import (
	"context"
	"fmt"
	"io/ioutil"
	tools "v1/tools"

	dns "google.golang.org/api/dns/v1beta2"
	"google.golang.org/api/option"
)

func CreateClientDnsName(EnvType, DnsName, Ip string) (string, error) {

	context := context.Background()

	data, err := ioutil.ReadFile(tools.JsonCredentialUrl)

	if err != nil {
		return "", err
	}

	client, err := dns.NewService(context, option.WithCredentialsJSON(data))

	if err != nil {
		return "", err
	}

	dnsChangeRequest := &dns.Change{

		Additions: []*dns.ResourceRecordSet{

			{
				Name: EnvType + "-" + DnsName + tools.DNSNAME,
				Type: "A",
				Ttl:  int64(300),
				Rrdatas: []string{
					Ip,
				},
			},
		},
	}

	dnsclient, err := client.Changes.Create(tools.MONKEYSPROJECT, tools.DNSCLIENTZONE, dnsChangeRequest).Context(context).Do()

	if err != nil {
		return "", err
	}

	fmt.Println(dnsclient)

	return dnsclient.Status, err
}

/* func CreateRecordSetHandler(w http.ResponseWriter, r *http.Request) {

	var K8sData tools.K8sData

	var response tools.Response

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}
	if body != nil {
		if err := json.Unmarshal(body, &K8sData); err != nil {
			//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		status, err := CreateClientDnsName(K8sData.EnvType, K8sData.DnsName, K8)

		if err != nil {
			response.Message = "Record fail to create"
			response.Result = "Error"
			w.WriteHeader(http.StatusBadRequest)
			encodeData, _ := json.Marshal(response)
			fmt.Fprintf(w, string(encodeData))
			return
		}

		response.Message = "Record set created succesfully"
		response.Result = "Success"
		response.Error = status
		w.WriteHeader(http.StatusOK)
		encodeData, _ := json.Marshal(response)
		fmt.Fprintf(w, string(encodeData))
		return
	}

} */
