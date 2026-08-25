基于 Go 实现的 LabVent 项目，一款后端服务，完成实验室通风柜与危化品的联动监控、报警与审计管理。

LabVent 是实验室通风柜与危化品监控平台：视窗位移先落盘再联动排风，排风机主备
切换保持柜内负压，多柜联排按序启动，风阀按行程顺序动作，危化品效期与分级按规则
校验，气体报警按闩锁复归，房间压差按改造后基线判定，全部操作写入审计。

## 构建与运行

```bash
go build -mod=vendor ./...
go test -mod=vendor -count=1 ./...
go vet -mod=vendor ./...
```

启动服务：

```bash
go run ./cmd/labvent -addr :18056 -data ./data
```

健康检查：`curl http://localhost:18056/healthz`

## Docker

```bash
bash build_benzhi_docker.sh labvent linux/amd64
docker run --rm -p 18056:18056 labvent bash -c 'go run ./cmd/labvent -addr :18056 -data /tmp/data'
```
