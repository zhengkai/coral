# 各种想法

淘汰数量：满了后进一个出一个，还是超过一定数量（如 1%）再批量清

过期清理频率：清理的最小时间间隔

过期清理程序在清完当前后，可得出最近下次清理时间，wait 到那个时间再执行

load：抓取上游的并发数有限制

统计：可选项