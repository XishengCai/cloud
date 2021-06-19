package main

import (
	"flag"
	env "github.com/alibabacloud-go/darabonba-env/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v2/client"
	console "github.com/alibabacloud-go/tea-console/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
)

var (
	regionId     = flag.String("regionId", "localhost:6379", "redis hostport")
	instanceName = flag.String("instanceName", "0", "redis database")
	password     = flag.String("password", "Root@123", "ecs password")
	instanceType = flag.String("instanceType", ":5040", "hostport to listen for HTTP JSON API")
	zoneId       = flag.String("zoneId", ":5040", "hostport to listen for HTTP JSON API")
	imageId      = flag.String("imageId", ":5040", "hostport to listen for HTTP JSON API")
	dryRun       = flag.Bool("dryRun", false, "")
)

func Initialization(accessKeyId *string, accessKeySecret *string, regionId *string) (_result *ecs20140526.Client, _err error) {
	config := &openapi.Config{}
	// 您的AccessKey ID
	config.AccessKeyId = accessKeyId
	// 您的AccessKey Secret
	config.AccessKeySecret = accessKeySecret
	// 您的可用区ID
	config.RegionId = regionId
	_result = &ecs20140526.Client{}
	_result, _err = ecs20140526.NewClient(config)
	return _result, _err
}

func CreateInstanceSample(client *ecs20140526.Client, regionId *string, name *string, password *string, instanceType *string, zoneId *string, imageId *string, dryRun *bool) (_err error) {
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		request := &ecs20140526.CreateInstanceRequest{
			InstanceName: name,
			Password:     password,
			InstanceType: instanceType,
			RegionId:     regionId,
			ZoneId:       zoneId,
			ImageId:      imageId,
			DryRun:       dryRun,
		}
		response, _err := client.CreateInstance(request)
		if _err != nil {
			return _err
		}

		console.Log(tea.String("创建实例成功，实例ID:" + tea.StringValue(response.Body.InstanceId)))

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.SetErrMsg(tryErr.Error())
		}
		console.Log(error.Message)
	}
	return _err
}

func DescribeInstancesSample(client *ecs20140526.Client, regionId *string) (_err error) {
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		request := &ecs20140526.DescribeInstancesRequest{
			RegionId: regionId,
		}
		response, _err := client.DescribeInstances(request)
		if _err != nil {
			return _err
		}

		instances := response.Body.Instances.Instance
		for _, instance := range instances {
			console.Log(tea.String("实例ID:" + tea.StringValue(instance.InstanceId) + "信息"))
			console.Log(tea.String("  状态:" + tea.StringValue(instance.Status)))
			console.Log(tea.String("  CPU:" + tea.ToString(tea.Int32Value(instance.Cpu))))
			console.Log(tea.String("  内存:" + tea.ToString(tea.Int32Value(instance.Memory)) + "MB"))
			console.Log(tea.String("  规格:" + tea.StringValue(instance.InstanceType)))
			console.Log(tea.String("  系统:" + tea.StringValue(instance.OSType) + "(" + tea.StringValue(instance.OSName) + ")"))
		}

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.SetErrMsg(tryErr.Error())
		}
		console.Log(error.Message)
	}
	return _err
}

func create() error {
	client, _err := Initialization(env.GetEnv(tea.String("ACCESS_KEY_ID")), env.GetEnv(tea.String("ACCESS_KEY_SECRET")), regionId)
	if _err != nil {
		return _err
	}

	// 1.创建实例
	console.Log(tea.String("--------------------创建实例--------------------"))
	_err = CreateInstanceSample(client, regionId, instanceName, password, instanceType, zoneId, imageId, dryRun)
	if _err != nil {
		return _err
	}
	// 2.等待实例创建成功
	_err = util.Sleep(tea.Int(1000))
	if _err != nil {
		return _err
	}
	// 2.查询实例列表
	console.Log(tea.String("--------------------查询实例列表--------------------"))
	_err = DescribeInstancesSample(client, regionId)
	if _err != nil {
		return _err
	}
	return _err
}

func main() {
	err := create()
	if err != nil {
		panic(err)
	}
}
