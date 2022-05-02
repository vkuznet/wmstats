module github.com/vkuznet/wmstats

go 1.18

require (
	github.com/dmwm/cmsauth v0.0.0-20220120183156-5495692d4ca7
	github.com/fatih/set v0.2.1
	github.com/gorilla/mux v1.8.0
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/prometheus/procfs v0.7.3
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/ulule/limiter/v3 v3.10.0
	github.com/vkuznet/http-logging v0.0.0-20210729230351-fc50acd79868
)

require (
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/lestrrat-go/strftime v1.0.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/tklauser/numcpus v0.4.0 // indirect
	github.com/vkuznet/x509proxy v0.0.0-20210801171832-e47b94db99b6 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9 // indirect
)

replace github.com/ulule/limiter/v3 => github.com/vkuznet/limiter/v3 v3.10.2
