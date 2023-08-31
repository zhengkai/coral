zhengkai/coral
======

[![Unit testing](https://github.com/zhengkai/coral/actions/workflows/unit-testing.yml/badge.svg)](https://github.com/zhengkai/coral/actions/workflows/unit-testing.yml)
[![Coverage Status](https://coveralls.io/repos/github/zhengkai/coral/badge.svg?branch=v2)](https://coveralls.io/github/zhengkai/coral?branch=v2)
[![Go Reference](https://pkg.go.dev/badge/github.com/zhengkai/coral/v2.svg)](https://pkg.go.dev/github.com/zhengkai/coral/v2)

Coral 是一个简单的 golang cache 类，无任何依赖，主要关注并发问题，即同时有 n 个 get 同时请求一个 key 时，只请求一次上游，其他的 get 都等待上游返回。

```
go get github.com/zhengkai/coral/v2
```
