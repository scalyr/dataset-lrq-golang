# DataSet LRQ client

This Golang package implements a simple client for the DataSet LRQ api

## Examples

### Log request

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
        fmt.Printf("%+v\n", log)
}
```

### Paginated log request

```golang
limit := 10

var cursor *string
for {
        var logs []lrq.LogResponseMatch
        var err error
        logs, cursor, err = client.DoLogRequestPaginated(ctx, lrq.LogRequestAttribs{
                Filter:    &filter,
                StartTime: stringToTime("2022-11-02T12:45:00-04:00"),
                EndTime:   stringToTime("2022-11-02T13:00:00-04:00"),
                Limit:     &limit,
        }, cursor)
        if err != nil {
                panic(err)
        }

        for _, log := range logs {
                fmt.Printf("%+v\n", log)
        }

        if len(logs) == 0 {
                break
        }
}
```

### Top facets request

```golang
numFacets := 5
valsPerFacet:= 3

facets, err := client.DoTopFacetsRequest(ctx, lrq.TopFacetsRequestAttribs{
        StartTime:         stringToTime("2022-11-02T12:45:00-04:00"),
        EndTime:           stringToTime("2022-11-02T13:00:00-04:00"),
        NumFacets:         &numFacets,
        NumValuesPerFacet: &valsPerFacet,
})
if err != nil {
        panic(err)
}

for _, facet := range facets {
        fmt.Printf("%+v\n", facet)
}
```

### Facet values request

```golang
values, err := client.DoFacetValuesRequest(ctx, "session", lrq.FacetValuesRequestAttribs{
        StartTime: stringToTime("2022-11-02T12:45:00-04:00"),
        EndTime:   stringToTime("2022-11-02T13:00:00-04:00"),
})
if err != nil {
        panic(err)
}

for _, value := range values {
        fmt.Printf("%+v\n", value)
}
```

### Plot request

```golang
slices := 10
breakdownFacet := "session"

plotData, err := client.DoPlotRequest(ctx, "count(tag='audit')", lrq.PlotRequestAttribs{
        StartTime:      stringToTime("2022-11-02T12:45:00-04:00"),
        EndTime:        stringToTime("2022-11-02T13:00:00-04:00"),
        Slices:         &slices,
        BreakdownFacet: &session,
})
if err != nil {
        panic(err)
}

fmt.Printf("%+v\n", plotData.XAxis)
for _, plot := range plotData.Plots {
        fmt.Printf("%+v\n", plot)
}
```

### PowerQuery plot request

```golang
pq := "tag='audit' | group count=count() by timestamp=timebucket('1m')"
plotData, err := client.DoPQPlotRequest(ctx, pq, lrq.PQRequestAttribs{
        StartTime: stringToTime("2022-11-02T12:45:00-04:00"),
        EndTime:   stringToTime("2022-11-02T13:00:00-04:00"),
})
if err != nil {
        panic(err)
}

fmt.Printf("%+v\n", plotData.XAxis)
for _, plot := range plotData.Plots {
        fmt.Printf("%+v\n", plot)
}
```

### PowerQuery table request

```golang
pq := "tag='audit' | group count=count() by timestamp=timebucket('1m')"
data, err := client.DoPQTableRequest(ctx, pq, lrq.PQRequestAttribs{
        StartTime: stringToTime("2022-11-02T12:45:00-04:00"),
        EndTime:   stringToTime("2022-11-02T13:00:00-04:00"),
})
if err != nil {
        panic(err)
}

fmt.Printf("%+v\n", data.Columns)
fmt.Printf("%+v\n", data.Values)
if len(data.Warnings) > 0 {
        fmt.Printf("%+v\n", data.Warnings)
}
```

### Distribution request

```golang
data, err := client.DoDistRequest(ctx, lrq.DistRequestAttribs{
        StartTime: stringToTime("2022-11-02T12:45:00-04:00"),
        EndTime:   stringToTime("2022-11-02T13:00:00-04:00"),
})
if err != nil {
        panic(err)
}

fmt.Printf("%+v\n", data)
```
