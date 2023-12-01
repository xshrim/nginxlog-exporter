# nginxlog-exporter

## 概述

nginxlog-exporter将nginx的access日志(或类似日志格式)转换为prometheus指标，从而实现对nginx所代理的后端服务的接口请求指标的无侵入式采集

nginxlog-exporter基于开源工具[martin-helmich/prometheus-nginxlog-exporter](https://github.com/martin-helmich/prometheus-nginxlog-exporter)，其功能特性、使用方法和配置文件与该工具几乎相同，区别仅在于nginxlog-exporter新增了对nginx日志的过滤功能域名证书有效期指标

## 编译

```shell
// 交叉编译
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o nginxlog-exporter main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o nginxlog-exporter.exe main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o nginxlog-exporter main.go
```

## ## 运行

1. 二进制执行

   ```bash
   chmod a+x nginxlog-exporter
   ./nginxlog-exporter -config-file nginxlog-exporter.yaml
   ```
2. 作为服务执行

   ```bash
   cat > /etc/systemd/system/nginxlog-exporter.service <<EOF
   [Unit]
   Description=NGINX metrics exporter for Prometheus
   After=network-online.target

   [Service]
   ExecStart=/usr/local/bin/nginxlog-exporter -config-file /etc/nginxlog-exporter.yaml
   ExecReload=/bin/kill -HUP $MAINPID
   Restart=always
   ProtectSystem=no
   CapabilityBoundingSet=

   [Install]
   WantedBy=multi-user.target
   EOF

   chmod a+x nginxlog-exporter
   mv nginxlog-exporter /usr/local/bin/
   mv nginxlog-exporter.yaml /etc/nginxlog-exporter.yaml
   mv nginxlog-exporter.service /etc/systemd/system
   systemctl start nginxlog-exporter
   ```

## 用法

详细使用方法请参考[martin-helmich/prometheus-nginxlog-exporter](https://github.com/martin-helmich/prometheus-nginxlog-exporter)的说明

此处仅对日志过滤功能的使用进行说明，日志过滤仅需在其配置文件中加入**filter_configs**字段（hcl格式则是filter字段）

```yaml
listen: # exporter监听配置(指标输出)
  port: 4040 # 服务端口
  address: "0.0.0.0" # 服务地址
  metrics_endpoint: "/metrics" # 指标输出路径

consul: # 服务发现，用于将此exporter注册到consul
  enable: false # 无需启用
  address: "localhost:8500" # consul服务发现server地址
  datacenter: dc1
  scheme: http
  token: ""
  service: # 注册自身
    id: "nginx-exporter"
    name: "nginx-exporter"
    address: "192.168.3.1"
    tags: ["foo", "bar"]

namespaces: # 以命名空间表示指标划分(可以团队为维度，也可以业务系统为维度)
  - name: demo_backend # 命名空间名称为指标名称前缀
    domains:    # 域名证书检测
      - demo.example.com:8443
    format: '$remote_addr $remote_user "$http_uname" "$http_uid" [$time_iso8601] $status "$request" $http_host $request_time "$upstream_response_time" $body_bytes_sent "$upstream_addr" "$upstream_status" "$http_referer" "$http_user_agent" "$http_x_forwarded_for"' # nginx access日志格式(与nginx配置文件中保持一致)
    filter_configs: # 将符合任一规则的日志行丢弃
      - from: line # 指定需要进行匹配的字段, line表示匹配整行
        match: ".*jndi:.*" # 正则匹配规则
      - from: request
        split: 2 # 对该nginx日志字段进行分隔，取第几个元素
        seperator: " " # nginx日志字段分隔符
        match: "^/.*(\\.js|\\.css|\\.ico|\\.rar|\\.zip|\\.tar|\\.tar\\.gz|\\.tgz|\\.gzip)" # 正则匹配规则
    relabel_configs: # 根据nginx日志字段为指标添加label
      - target_label: host # 指标label名称
        from: http_host # 该指标label基于哪个nginx日志字段
      - target_label: upstream
        from: upstream_addr
        # whitelist: ["127.0.0.1"]  # 字段白名单，字段值不在白名单内则为other
      - target_label: request
        from: request
        split: 2 # 对该nginx日志字段进行分隔，取第几个元素
        seperator: " " # nginx日志字段分隔符
        matches: # 对nginx日志字段指定元素的值进行过滤，符合条件的才会加入指标中
          - regexp: "^/(.*)/[0-9]+/?.*" # 正则匹配，不匹配则继续匹配下一条规则
            replacement: "/$1/:id" # 对匹配的字段元素值进行正则替换
          - regexp: "^/([^\\?\\n]*).*"
            replacement: "/$1"

    print_log: false # exporter打印debug日志

    source: # exporter采集的nginx日志文件
      files: # 支持file和syslog两种模式
        - /var/log/nginx/access.log
        # - ./test.log
    metrics_override: # 指标名覆盖(用于多个命名空间使用相同的指标名)
      prefix: "demolog" # 将指标名中的namespace前缀替换
    namespace_label: "namespace" # 将命名空间名称作为指标label
    labels: # 额外为命名空间下的所有指标添加label
      biz: "trade"
      sys: "demo"
      svc: "backend"
    histogram_buckets: [.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 20] # histogram类型指标的bucket区间设计

```

> 因过滤或解析失败而丢弃的日志行数均会统计在**DiscardTotal**指标中，并以*reason标签进行区分*
