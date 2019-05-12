# line-pay-sdk-go

line-pay-sdk-go is a Go client library for accessing the [LINE Pay API](https://pay.line.me/jp/developers/documentation/download/tech).

## Usage

```go
import "github.com/gotokatsuya/line-pay-sdk-go/linepay"

func main() {
    pay, err := linepay.New("<channel id>", "<channel secret>")
    ...
}
```

## License

This library is distributed under the MIT license.
