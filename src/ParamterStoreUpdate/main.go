package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const constLayer = "/Cloud/prod/ConstantLayerVersionArn"
const layerVer = "/Cloud/prod/LAYERVERSIONARN"

func main() {
	params, err := getParameterStore()
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(params)
	ServiceList := getServiceParameter(params)
	serviceKeys := make([]string, 0, len(ServiceList))
	for k := range ServiceList {
		serviceKeys = append(serviceKeys, k)
	}
	sort.Strings(serviceKeys)
	// fmt.Println(ServiceList)
	ite := 0
	changeList := make(map[int][]string)
	for {
		for _, key := range serviceKeys {
			fmt.Printf("Service: %s\n", key)
		}
		fmt.Printf("更新するサービスを選択してください\n")
		fmt.Printf("パラメータの更新を実行、もしくはCLIを終了するにはExitを入力\n")
		fmt.Printf("更新するサービス: ")
		var i string
		fmt.Scan(&i)
		if i == "Exit" {
			break
		}
		fmt.Printf("更新するパラメータを選択してください\n")
		if _, ok := ServiceList[i]; ok {
			for _, list := range ServiceList[i] {
				fmt.Printf("%s", list)
			}
		} else {
			continue
		}

		fmt.Printf("はじめに戻るにはBackを入力\n")
		fmt.Printf("更新するサービス: ")
		var j string
		fmt.Scan(&j)
		if j == "Back" {
			continue
		}

		fmt.Printf("現在のパラメータ値はこちら、変更したい値を入力してください\n")
		fmt.Printf("はじめに戻るにはBackを入力\n")
		_, err := getParameterValue(i, j)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		var k string
		fmt.Printf("入力値: ")
		fmt.Scan(&k)
		if k == "Back" {
			continue
		}
		changeList[ite] = append(changeList[ite], i)
		changeList[ite] = append(changeList[ite], j)
		changeList[ite] = append(changeList[ite], k)
		ite++
	}
	// fmt.Println(changeList)
	for _, val := range changeList {
		fmt.Printf("サービス「%v」の「%v」が「%v」に更新されます\n", val[0], val[1], val[2])
	}
	fmt.Printf("パラメータの更新をしますか (y/n): ")
	var l string
	fmt.Scan(&l)
	if l == "y" {
		changeParameterStore(changeList)
	}
}

func getParameterStore() ([]string, error) {
	queryPath := "/Cloud/prod/"

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ssm.New(sess)

	paramList := &ssm.GetParametersByPathInput{
		MaxResults: aws.Int64(10),
		Path:       aws.String(queryPath),
		Recursive:  aws.Bool(true),
	}

	var params []string
	for {
		result, err := svc.GetParametersByPath(paramList)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		for i := 0; i < len(result.Parameters); i++ {
			if strings.Contains(strings.SplitN(strings.SplitN(result.Parameters[i].GoString(), ",", 6)[3], ":", 2)[1], constLayer) || strings.Contains(strings.SplitN(strings.SplitN(result.Parameters[i].GoString(), ",", 6)[3], ":", 2)[1], layerVer) {
				continue
			}
			params = append(params, strings.SplitN(strings.SplitN(result.Parameters[i].GoString(), ",", 6)[3], ":", 2)[1]+"\n")
		}
		if result.NextToken == nil {
			break
		}
		paramList.SetNextToken(*result.NextToken)

	}
	// fmt.Println(params)
	return params, nil
}

func getServiceParameter(parameterList []string) map[string][]string {
	a := make(map[string][]string)
	var Service string
	for i := 0; i < len(parameterList); i++ {
		Service = strings.SplitN(parameterList[i], "/", 5)[3]
		a[Service] = append(a[Service], strings.Replace((strings.SplitN(parameterList[i], "/", 5)[4]), "\"", "", 1))
	}
	// fmt.Println(a[Service])
	return a
}

func getParameterValue(i, j string) (string, error) {
	queryPath := "/Cloud/prod/" + i + "/" + j

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ssm.New(sess)

	queryParam := &ssm.GetParameterInput{
		Name: aws.String(queryPath),
	}
	result, err := svc.GetParameter(queryParam)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	fmt.Println(*result.Parameter.Value)
	return "", nil
}

func changeParameterStore(changeList map[int][]string) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ssm.New(sess)
	for i := 0; i < len(changeList); i++ {
		putParameterDetail := &ssm.PutParameterInput{
			Name:      aws.String("/Cloud/prod/" + changeList[i][0] + "/" + changeList[i][1]),
			Overwrite: aws.Bool(true),
			Value:     aws.String(changeList[i][2]),
			Type:      aws.String("String"),
		}
		_, err := svc.PutParameter(putParameterDetail)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
