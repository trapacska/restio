# RESTIO

If you need to optimize writing Go http handlers. It is a wrapper around built-in `http.Handler` what means that this package supports all 3rd party mux packages as well. (for example gorilla/mux)

___

## Return handling

The `error` return value will be used to automatically set statuscode and write JSON error message with the error into the response body. If your handler has any `interface{}` return value set then it will be Encoded to JSON and written into the response body. 

You can use both at the same time. Accepted return values:
```
func handler(w http.ResponseWriter, r *http.Request)
func handler(w http.ResponseWriter, r *http.Request) error
func handler(w http.ResponseWriter, r *http.Request) interface{}
func handler(w http.ResponseWriter, r *http.Request) (interface{}, error)
```

### Example
```
func handler(w http.ResponseWriter, r *http.Request) error {
    return fmt.Errorf("oops... an error")
}
```

> Status: 501 InternalServerError
> ```
> {
>     "error": "oops... an error"
> }
> ```

--- 
```
type Result struct {
    Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) Result {
    return Result{Message: "WOW!"}
}
```

> Status: 200 Success
> ```
> {
>     "message": "WOW!"
> }
> ```

## Arguments handling

By default you have 2 arguments but if you add your custom `interface{}` argument as first, it will be populated by the JSON data sent in the request body.

Accepted arguments:
```
func handler(w http.ResponseWriter, r *http.Request)
func handler(i interface{}, w http.ResponseWriter, r *http.Request)
```

### Example
```
type Result struct {
    Message string `json:"message"`
}

func handler(res Result, w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, res.Message)
}
```

Request body data:
```
{
    "message": "WOW!"
}
```

> Status: 200 Success
> ```
> WOW!
> ```

## Usage

Wrap your existing handlers in `restio.Endpoint(..handler..)` and that's all.

### Example

```
r := mux.NewRouter()

r.Handle("/test", restio.Endpoint(testHandler)).Methods("GET")
```
