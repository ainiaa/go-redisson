package redis

import (
	"context"
	"time"

	gredis "github.com/go-redis/redis/v8"
	"github.com/hashicorp/go-uuid"
)

const (
	LockScript      = "if redis.call('exists', KEYS[1]) == 0  then return redis.call('setex', KEYS[1], unpack(ARGV)) else return '-1' end"
	UnlockScript    = "if redis.call('get', KEYS[1]) == ARGV[1]  then return redis.call('del', KEYS[1]) or true end"
	RenewLockScript = "if redis.call('get', KEYS[1]) == ARGV[1]  then return redis.call('expire', KEYS[1],ARGV[2]) or true end"
	LockPrefix      = "lock:"
)

var (
	lockScripter           *gredis.Script
	unlockScripter         *gredis.Script
	renewScripter          *gredis.Script
)

/**
 * 执行加锁操作
 *
 * param: []string args
 * param: int64    lockTime
 * param: string    identifier
 * return: string, error
 */
func (r *Redis) doLock(ctx context.Context, args []string, lockTime int64, identifier string) (res string, err error) {
	var result interface{}
	result, err = r.GetLockScripter(ctx).Run(ctx, r, args, lockTime, identifier).Result()

	if err != nil {
		return "", err
	}
	res = result.(string)

	if res == "-1" { // 若已存在,直接返回
		return "", ErrExitsLock
	}

	if res == "OK" { //获取锁成功
		return identifier, nil
	}

	return "", ErrAcquiredLock
}

/**
 * 简单获取锁
 *
 * param: string lockName
 * param: int    lockTime
 * return: string
 * return: error
 */
func (r *Redis) LockSingle(ctx context.Context, lockName string, lockTime int64) (identifier string, err error) {
	identifier, _ = uuid.GenerateUUID()
	args := []string{LockPrefix + lockName}
	return r.doLock(ctx, args, lockTime, identifier)
}

/**
 * 获取锁
 *
 * param: string lockName
 * param: int    lockTime
 * param: int    acquireTime
 * return: string
 */
func (r *Redis) Lock(ctx context.Context, lockName string, lockTime int64, acquireTime int) (string, error) {
	identifier, _ := uuid.GenerateUUID()
	return r.LockWithId(ctx, lockName, identifier, lockTime, acquireTime)
}

/**
 * 获取锁
 *
 * param: string lockName
 * param: int    lockTime
 * param: int    acquireTime
 * return: string
 */
func (r *Redis) LockWithId(ctx context.Context, lockName string, identifier string, lockTime int64, acquireTime int) (string, error) {
	if acquireTime == 0 {
		return r.LockSingle(ctx, lockName, lockTime)
	}
	args := []string{LockPrefix + lockName}
	_, err := r.doLock(ctx, args, lockTime, identifier)
	if err == nil {
		return identifier, err
	}

	ticker := time.NewTicker(time.Duration(20) * time.Millisecond)
	defer func() {
		ticker.Stop()
	}()
	timer := time.After(time.Duration(acquireTime) * time.Second)
	for {
		select {
		case <-timer: //超时
			return "", ErrAcquiredLockTimeout
		case <-ticker.C: //执行时间执行一次
			_, err := r.doLock(ctx, args, lockTime, identifier)
			if err == nil {
				return identifier, err
			}
		}
	}
}

/**
 * 释放锁
 *
 * param: string lockName
 * param: string lockId
 * return: error
 */
func (r *Redis) Unlock(ctx context.Context, lockName, lockId string) (err error) {
	args := []string{LockPrefix + lockName}
	_, err = r.GetUnlockScripter(ctx).Run(ctx, r, args, lockId).Result() // 删除成功 返回 1, nil 删除失败 nil, gredis.Nil

	if err == gredis.Nil { //没有找到对应的锁或者锁和给定的lockId不匹配
		err = nil
	}

	return err
}

/**
 * 延长锁
 *
 * param: string lockName
 * param: string lockId
 * param: int    renameTime
 * return: error
 */
func (r *Redis) RenewLock(ctx context.Context, lockName, lockId string, renameTime int) (err error) {
	args := []string{LockPrefix + lockName}
	_, err = r.GetRenewScripter(ctx).Run(ctx, r, args, lockId, renameTime).Result()

	return err
}

func (r *Redis) LoadLockScript(ctx context.Context) (*gredis.Script, error) {
	if lockScripter == nil {
		lockScripter = gredis.NewScript(LockScript)
	}
	_, err := lockScripter.Load(ctx, r).Result()
	if err != nil {
		return nil, err
	}

	return lockScripter, nil
}

func (r *Redis) LoadUnLockScript(ctx context.Context) (*gredis.Script, error) {
	if unlockScripter == nil {
		unlockScripter = gredis.NewScript(UnlockScript)
	}
	_, err := unlockScripter.Load(ctx, r).Result()
	if err != nil {
		return nil, err
	}
	return unlockScripter, nil
}

func (r *Redis) LoadRenewScript(ctx context.Context) (*gredis.Script, error) {
	if renewScripter == nil {
		renewScripter = gredis.NewScript(RenewLockScript)
	}
	_, err := unlockScripter.Load(ctx, r).Result()
	if err != nil {
		return nil, err
	}
	return renewScripter, nil
}

func (r *Redis) GetLockScripter(ctx context.Context) *gredis.Script {
	var err error
	if lockScripter == nil {
		lockScripter, err = r.LoadLockScript(ctx)
		if err != nil {
			panic(err)
		}
	}
	return lockScripter
}

func (r *Redis) GetUnlockScripter(ctx context.Context) *gredis.Script {
	var err error
	if unlockScripter == nil {
		unlockScripter, err = r.LoadUnLockScript(ctx)
		if err != nil {
			panic(err)
		}
	}
	return unlockScripter
}

func (r *Redis) GetRenewScripter(ctx context.Context) *gredis.Script {
	var err error
	if renewScripter == nil {
		renewScripter, err = r.LoadRenewScript(ctx)
		if err != nil {
			panic(err)
		}
	}
	return renewScripter
}
