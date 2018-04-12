![](./logo.jpg)

# AssassinGo

AssassinGo is an extensible and concurrency information gathering and vulnerability scanning framework.

Just for learn.

## Features

- [x] Retrieve Server and Powered by Headers
- [x] Detect CMS Version
- [ ] Subdomain Scan
- [ ] Google Map
- [ ] Dir Scan
- [ ] Whois Lookup
- [x] Port Scan
- [x] Crawl the Paramed URLs
- [x] Basic SQLi Check
- [x] Basic XSS Check
- [x] Google Hacking
- [x] Brute Force
- [x] POC Interface
- [x] Web GUI(using WebSocket)

## API

### AJAX

Path | Method | Func | Params | Return
----- | ----- | ----- | ----- | -----
/api/target | POST | set a target | target=xxx / target=t1,t2... | nil
/api/poc | GET | get poc list | nil | {"pocs": pocs}
/api/poc | POST | set poc | poc=xxx | nil

### WebSocket

Path | Func | Params | Return
----- | ----- | ----- | -----
/ws/info/basic | get ip and webserver | nil | {"ip": ip, "webserver": webserver}
/ws/info/cms | detect cms | nil | {"cms": cms}
/ws/info/port | ports scan | {method:"tcp"/"syn"} | {"ports": ports}
/ws/info/crawl | crawl paramed urls | {max_depth: 4} | {"url": url}
/ws/vul/sqli | check sqli | nil | {"url": url}
/ws/vul/xss | check xss | nil | {"url": url}
/ws/intrude | brute force | {header:"xxx", payload:"p1,p2...", gort_count="10"} | {"payload": payload, "resp_status": status, "resp_len": len}
/ws/seek | seek targets | {query:"", se:"bing/google", max_page: 10} | {"urls": urls}
/ws/poc/run | run poc | nil | {"host": host}

## State

developing...

## Demo

![](./demo.png)