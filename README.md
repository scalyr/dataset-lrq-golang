# DataSet LRQ client

This Golang package implements a simple client for the DataSet LRQ api

## Examples

Here is a (non-paginated) log request example

```golang
# Client initialization using a log access api key
client := lrq.NewClient("https://app.scalyr.com", "<apikey>")

stringToTime := func(s string) *time.Time {
        t, err := time.Parse(time.RFC3339, s)
        if err != nil {
                panic(err)
        }
        return &t
}

# The provided context allows request cancellation and/or timeout
ctx := context.Background()

filter := "tag='audit'"

logs, err := client.DoLogRequest(ctx, lrq.LogRequestAttribs{
        Filter:    &filter,
        StartTime: stringToTime("2022-11-02T12:45:00-04:00"),
        EndTime:   stringToTime("2022-11-02T13:00:00-04:00"),
})
if err != nil {
        panic(err)
}

for _, log := range logs {
        fmt.Printf("%v\n", log)
}
```
