<h1 align="center">HedgeDamnRequest</h1>
<br/>

_Using a hedged HTTP client, it sends requests with a configurable delay between attempts, logs detailed request and response information, and calculates the speed of each server._

## Technologies/Library used
* Golang
* [hedgedhttp](https://github.com/cristalhq/hedgedhttp)

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





