![](./logo.jpg)

# AssassinGo

AssassinGo is an extensible and concurrency information gathering and vulnerability scanning framework.

Just for learn.

## Features

- [x] Retrieve Server and Powered by Headers
- [x] Detect CMS Version
- [x] Port Scan
- [x] Trace Route and Mark on Google Map
- [ ] Subdomain Scan
- [x] Dir Scan
- [ ] Site Map
- [x] Whois Lookup
- [x] Crawl the Paramed URLs
- [x] Basic SQLi Check
- [x] Basic XSS Check
- [x] Intruder
- [x] Seek Targets by Google
- [x] POC Interface
- [x] Web GUI(using WebSocket)

## API

### AJAX

Path | Method | Func | Params | Return
----- | ----- | ----- | ----- | -----
/api/target | POST | set a target | target=xxx / target=t1,t2... | nil
/api/poc | GET | get poc list | nil | {"pocs": ["seacmsv654", "struct056"...]}
/api/poc | POST | set poc | poc=xxx | nil

### WebSocket

Path | Func | Params | Return
----- | ----- | ----- | -----
/ws/info/basic | get ip and webserver | nil | {"ip": "192.168.1.1", "webserver": "nginx"}
/ws/info/whois | whois | nil | {"domain":"example.com","registrar_name":"alibaba", "admin_name":"xiaoming", "admin_email":"a@qq.com", "admin_phone":"+86.12312345678", "created_date":"2016-07-28T12:57:53.0Z","expiration_date":"2018-07-28T12:57:53.0Z", "ns":"dns9.hichina.com", "state":"clienttransferprohibited"}
/ws/info/cms | detect cms | nil | {"cms": "wordpress"}
/ws/info/port | ports scan | {method:"tcp"/"syn"} | {"port": "80", "service": "http"}
/ws/info/tracert | trace route and mark on google map | nil | {"ttl": 1, "addr": 192.168.1.1, "elapsed_time": 22720440, "country": China, "lat": 34.2583,"long": 116.1614}
/ws/info/dirb | brute force dir | {"gort_count":20, "dict":"php"}; {"stop":1} | {"path": "admin.php", "resp_status": 200, "resp_len": 110}
/ws/info/crawl | crawl paramed urls | {"max_depth": 4} | {"url": "example.com/?id=1"}
/ws/vul/sqli | check sqli | nil | {"url": "example.com/?id=1}
/ws/vul/xss | check xss | nil | {"url": "example.com/?id=1}
/ws/intrude | brute force | {"header": "GET / HTTP/1.1 ...", "payload": "p1,p2...", "gort_count": "10"}; {"stop":1}| {"payload": 1, "resp_status": 200, "resp_len": 110}
/ws/seek | seek targets | {"query": "biu", "se": "bing/google", "max_page": 10} | {"urls": urls}
/ws/poc/run | run poc | nil | {"host": "example.com"}

## State

developing...

## Demo

![](./demo.png)