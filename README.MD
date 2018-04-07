![](./logo.jpg)

# AssassinGo

AssassinGo is an extensible and concurrency information gathering and vulnerability scanning framework.

Just for learn.

## Features

- [x] Retrieve Server and Powered by Headers
- [x] Detect CMS Version
- [ ] DNS Records Dump
- [ ] Whois Lookup
- [x] Port Scan
- [x] Crawl the Paramed URLs
- [x] Basic SQLi Scan
- [x] Basic XSS Scan
- [ ] Google Hacking
- [ ] Brute Force
- [x] Friendly API.
- [ ] POC Lib
- [ ] Web GUI(using WebSocket)

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
/ws/intrude/ | brute force | header=xxx&payload=xxx | {"payload": payload, "resp_status": status, "resp_len": len}

## State

developing...

## Demo

![](./demo.png)