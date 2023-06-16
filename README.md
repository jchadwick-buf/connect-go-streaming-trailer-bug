# Connect Go streaming trailer canonicalization bug
This repository demonstrates an issue with Connect Go wherein trailers parsed from Connect streaming are not canonicalized as expected.

If the bug is present, the following output will be seen when running `go run .`:
```go
Unary:
- Trailer().Get("lowercase"): test
- Trailer()["Lowercase"][0]: test
Streaming:
- Trailer().Get("lowercase"): 
- Trailer()["lowercase"][0]: test
```