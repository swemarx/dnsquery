all:
	go build dnsquery.go
	strip dnsquery
