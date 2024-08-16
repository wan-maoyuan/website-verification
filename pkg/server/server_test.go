package server

import "testing"

func TestEncodingUrl(t *testing.T) {
	urlStr := "http://47.239.194.150:31500/api/com.hello.jim?pid=jimobi_int&clickid=94655239-b0bb-4d58-900f-29118b7e62db&ua=Mozilla/5.0 (Linux; Android 14; V2323A Build/UP1A.231005.007; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/116.0.0.0 Mobile Safari/537.36&ip=224.36.155.11&android_id=3b9774a8-a6a6-4b7e-9875-925cfddbf830"
	encodingUrl := encodingUrl(urlStr)
	t.Logf(encodingUrl)
}
