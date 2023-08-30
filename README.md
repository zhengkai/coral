github.com/zhengkai/coral
======

Coral 是一个简单的 golang cache 类，主要关注并发问题，即同时有 n 个 get 同时请求一个 key 时，只请求一次上游，其他的 get 都等待上游返回。

```
go get github.com/zhengkai/coral/v2
```
