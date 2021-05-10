package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ainiaa/bytesconv"
	jsoniter "github.com/json-iterator/go"

	"github.com/ainiaa/go-redission/conf"
)

func TestRedis_LoadLockScript(t *testing.T) {
	var ctx = context.Background()
	c := getConf()
	InitOnceRedis(ctx, &c)
	_, err := GetRedis().LoadLockScript(ctx)
	if err != nil {
		t.Errorf("LoadLockScript error:%s", err.Error())
	} else {
		t.Log("LoadLockScript success")
	}
}

func TestRedis_Lock(t *testing.T) {
	var ctx = context.Background()
	c := getConf()
	InitOnceRedis(ctx, &c)
	type args struct {
		lockName    string
		lockTime    int64
		acquireTime int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test-lock", args{"test-lock", 100, 2}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRedis().Lock(ctx, tt.args.lockName, tt.args.lockTime, tt.args.acquireTime)
			err2 := GetRedis().Close()
			if err2 != nil {
				t.Errorf("close error = %v\n", err2)
			} else {
				t.Log("close success \n")
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Lock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Lock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_RenewLock(t *testing.T) {
	var ctx = context.Background()
	c := getConf()
	InitOnceRedis(ctx, &c)
	type args struct {
		lockName   string
		lockId     string
		renameTime int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetRedis().RenewLock(ctx, tt.args.lockName, tt.args.lockId, tt.args.renameTime); (err != nil) != tt.wantErr {
				t.Errorf("RenewLock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedis_Unlock(t *testing.T) {
	var ctx = context.Background()
	c := getConf()
	InitOnceRedis(ctx, &c)
	var err error
	var lockName1 = "abc"
	var lockTime1 = int64(10)
	var id1 string
	id1, err = GetRedis().LockSingle(ctx, lockName1, lockTime1)
	fmt.Printf("id1 GetLock error:%v\n", err)
	var lockName2 = "abcd"
	var lockTime2 = int64(10)
	var id2 string
	id2, err = GetRedis().LockSingle(ctx, lockName2, lockTime2)
	fmt.Printf("id2 GetLock error:%v\n", err)
	fmt.Printf("id1:%s\n", id1)
	fmt.Printf("id2:%s\n", id2)
	type args struct {
		lockName string
		lockId   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"unlock-id1", args{lockName1, id1}, false},
		{"unlock-id2", args{lockName2, "abcd"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetRedis().Unlock(ctx, tt.args.lockName, tt.args.lockId); (err != nil) != tt.wantErr {
				t.Errorf("Unlock() error = %v, wantErr %v", err, tt.wantErr)
			}
			GetRedis().Close()
		})
	}
}

func getConf() conf.Config {
	confStr := `
{
    "conn_type": 1,
    "alone": {
        "network": "alone",
        "addr": "127.0.0.1:6379",
        "username": "",
        "password": "123456",
	    "dial_timeout":10000
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

func TestRedis_Time(t *testing.T) {
	tiker := time.NewTicker(time.Second)
	for i := 0; i < 3; i++ {
		fmt.Println(<-tiker.C)
	}
}
