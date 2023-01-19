# beego-cache
The independent cache module from Beego.

So if you got any issues, please raise the issues on [Beego](https://github.com/beego/beego).

## TODO
- 分离单元测试和 e2e 测试
- redigo 和 redis-go
    - 做法1：直接用 redis-go 替换掉 redigo
    - 做法2：提供一个新的实现，redis-go 的实现：redis-go v7
- github action 要搞起来
- 删除 StartAndGC 方法，让用户直接创建对应的缓存实现 NewRedisCache(cmd redis.Cmdable, opts...Option),
- NewRedisCache(cfg string)
