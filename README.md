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

```
rdb := redis.NewClient(&redis.Options{
   Addr:     "192.168.56.105:6379",
   Password: "", // no password set
   DB:       0,  // use default DB

})
```
