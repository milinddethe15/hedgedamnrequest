<h1 align="center">HedgeDamnRequest</h1>
<br/>

_Using a hedged HTTP client, it sends requests with a configurable delay between attempts, logs detailed request and response information, and calculates the speed of each server._

## Technologies/Library used
* Golang
* [hedgedhttp](https://github.com/cristalhq/hedgedhttp)

## Hedged Requests

### What are Hedged requests?
Hedged requests are way to improve the speed and reliability of HTTP requests. When you send a hedged request, you send multiple identical requests to different servers or services, hoping that one of them will respond faster. The first response that comes back is used and the others are ignored or canceled.

### What it solves?
Hedged requests help solve the problem of unpredictable response times from servers. Sometimes, a server might be slow due to high traffic, network issues or other reasons. By sending multiple requests, it is more likely to get a fast response, reducing the overall waiting time.

### How to implement it?
- Set up a system where multiple identical requests are sent with a small delay between them.
- Use the first response that comes back and cancel the rest.

In this project, [hedgedhttp](https://github.com/cristalhq/hedgedhttp) is used to manage the hedged requests. You set a delay between each request (called the hedge delay) and specify how many requests to send (maximum attempts).

## Running locally
First, clone the repo:

```
git clone https://github.com/milinddethe15/hedgedamnrequest.git
```

Run locally:

```bash
go mod tidy
go run service/service.go
go run main.go #use another terminal tab
```

## Usage

Add list of mirror server links to config yaml file to find which mirror link is faster.

## Screenshot
![hedgedamnrequest](https://github.com/user-attachments/assets/a52e5445-f994-402f-b751-1df6566e3f2c)





