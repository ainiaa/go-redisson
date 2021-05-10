package main

import (
	"context"
	"fmt"

	"github.com/ainiaa/bytesconv"
	jsoniter "github.com/json-iterator/go"

	"github.com/ainiaa/go-redisson"
	"github.com/ainiaa/go-redisson/conf"
)

func main() {
	var err error
	var ctx = context.Background()
	c := getConf()
	err = redis.InitOnceRedis(ctx, &c)
	if err != nil {
		panic(fmt.Sprintf("InitOnceRedis error: %s", err))
	}
	_, err = redis.GetRedis().LoadLockScript(ctx)
	if err != nil {
		fmt.Printf("LoadLockScript error:%s\n", err.Error())
	} else {
		fmt.Printf("LoadLockScript success\n")
	}

	_, err = redis.GetRedis().LoadUnLockScript(ctx)
	if err != nil {
		fmt.Printf("LoadUnLockScript error:%s\n", err.Error())
	} else {
		fmt.Printf("LoadUnLockScript success\n")
	}

	err = redis.InitNamedRedis(ctx, "default", &c)
	if err != nil {
		panic(fmt.Sprintf("InitNamedRedis error: %s", err))
	}

	_, err = redis.GetNamedRedis("default").LoadLockScript(ctx)
	if err != nil {
		fmt.Printf("GetNamedRedis().LoadLockScript error:%s\n", err.Error())
	} else {
		fmt.Printf("GetNamedRedis().LoadLockScript success\n")
	}
	id, err := redis.GetNamedRedis("default").Lock(ctx, "abc", 10, 10)
	if err != nil {
		fmt.Printf("Lock:adc error:%s\n", err.Error())
	} else {
		fmt.Printf("Lock:adc success id:%s\n", id)
		err = redis.GetNamedRedis("default").Unlock(ctx, "abc", id)
		if err != nil {
			fmt.Printf("Unlock:adc id:%s error id:%s\n", id, err.Error())
		} else {
			fmt.Printf("Unlock:adc  id:%s success\n", id)
		}
	}

}

func getConf() conf.Config {
	confStr := `
{
    "conn_type": "alone",
    "alone": {
        "network": "",
        "addr": "127.0.0.1:6379",
        "username": "",
        "password": "123456"
    },
    "cluster": {
        "addrs": []
    }
}
`
	var c = conf.Config{}
	var json = jsoniter.ConfigDefault

	var err = json.Unmarshal(bytesconv.StringToBytes(confStr), &c)
	if err != nil {
		panic("config json.Unmarshal error:" + err.Error())
	}
	return c
}
