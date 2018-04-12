![](./logo.jpg)

# AssassinGo

AssassinGo is an extensible and concurrency information gathering and vulnerability scanning framework.

Just for learn.

## Features

- [x] Retrieve Server and Powered by Headers
- [x] Detect CMS Version
- [ ] Subdomain Scan
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
/api/target | POST | set a target | target=xxx | nil
/api/targets | POST | set targets | targets=t1,t2... | nil
/api/info/basic | GET | get ip and webserver | nil | {"ip": ip, "webserver": webserver}
/api/info/cms | GET | detect cms | nil | {"cms": cms}
/api/info/port | GET | ports scan | nil | {"ports": ports}
/api/poc | GET | get poc list | nil | {"pocs": pocs}
/api/poc | POST | set poc | poc=xxx | nil

### WebSocket

Path | Func | Params | Return
----- | ----- | ----- | -----
/ws/crawl | crawl paramed urls | nil | {"url": url}
/ws/vul/sqli | check sqli | nil | {"url": url}
/ws/vul/xss | check xss | nil | {"url": url}
/ws/intrude | brute force | {header:"xxx", payload:"p1,p2...", gort_count="10"} | {"payload": payload, "resp_status": status, "resp_len": len}
/ws/seek | seek targets | {query:"", se:"bing/google", max_page: 10} | {"urls": urls}
/ws/poc/run | run poc | nil | {"host": host}

## State

developing...

## Demo

![](./demo.png)