# Pool
*simplified worker pool*

[![GoDoc](https://godoc.org/github.com/selfquery/pool?status.png)](https://godoc.org/github.com/selfquery/pool)

# Install
`go get github.com/selfquery/pool`

# Example
```go
type test struct {
	Val int
}

func (t test) Process() {
	fmt.Println(t)
}

func (t test) Run(wait *sync.WaitGroup) {
	t.Process()
	wait.Done()
}
```
```go
func main() {
	p := pool.CreatePool([]pool.Details{
		test{1},
		test{1},
		test{1},
	}, 1)
	p.Run()
}

```

# Complex Example
```go 
type test struct {
    Output string
    Url string
    Attempts int
}

func (t test) Process() {

	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
	}

	res, err := client.Get(t.Url)
	if err != nil || res.StatusCode != 200 {
		if t.Attempts >= 3 {
			return
		} else {
			t.Attempts++
			t.Process()
			return
		}
	}
	defer res.Body.Close()

	file, err := os.Create(t.Output)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		panic(err)
	}

}

func (t test) Run(wait *sync.WaitGroup) {
	t.Process()
	wait.Done()
}
```