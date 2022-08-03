# distributed-lock
100行代码 基于redis 实现的分布式锁

功能接口：

| 方法名称   | 说明                                                     |
| ---------- | -------------------------------------------------------- |
| 1、NewLock | 新建锁对象。                                             |
| 2、Lock    | 加锁；如果为获取到锁，会阻塞。                           |
| 3、Unlock  | 解锁。                                                   |
| 4、Close   | 关闭，关闭后不能再次调用加锁解锁功能，只能再次重新申请。 |



使用方法：

1、创建redis对象

```go
rdb := redis.NewClient(&redis.Options{
   Addr:     "192.168.56.105:6379",
   Password: "", // no password set
   DB:       0,  // use default DB

})
```

2、创建锁对象

```go
rr, err := NewLock(rdb, "test-redis-lock")
if err != nil {
   t.Fatal(err)
}
```

```
func NewLock(rdb *redis.Client, key string) (r *redisLock, err error) 
参数说明：
rdb redis 客户端指针
key 一个锁的标识，标识不是最终的key值，实现时在后面有拼接字符标识。
```

3、获取锁

```go
r.Lock()
```

4、释放锁

```go
r.Unlock()
```

![](img\redis-lock-log.jpg)
