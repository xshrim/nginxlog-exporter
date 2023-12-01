package ssl

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

func Check(domain string) (string, string, string, string, string, string, float64, error) {
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	info := strings.Split(domain, ":")
	domain = info[0]
	port := "443"
	if len(info) > 1 {
		port = info[1]
	}

	client := &http.Client{
		Transport: &http.Transport{
			// 如果证书已过期，则仅在关闭证书校验的情况下链接才能建立成功
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true, MinVersion: tls.VersionTLS10},
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				addr = "127.0.0.1:" + port
				return (&net.Dialer{Timeout: 5 * time.Second, KeepAlive: 10 * time.Second}).DialContext(ctx, network, addr)
			},
		},
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("https://%s:%s", domain, port))
	if err != nil {
		return domain, port, "", "", "", "", 0, err
	}
	defer resp.Body.Close()

	for _, cert := range resp.TLS.PeerCertificates {
		// 检测服务器证书是否已经过期(CA证书过期时间会比服务器证书长)
		if !cert.IsCA {
			sn := fmt.Sprintf("%x", cert.SerialNumber)
			issuer := cert.Issuer.CommonName
			dnss := strings.Join(cert.DNSNames, ",")
			expire := cert.NotAfter.Local().Format("2006-01-02 15:04:05")
			remain := cert.NotAfter.Sub(time.Now()).Seconds()
			return domain, port, sn, issuer, dnss, expire, remain, cert.VerifyHostname(domain)
		}
	}

	return domain, port, "", "", "", "", 0, nil

}
