# snap-plugin-collector-skeleton
This plugin is meant to be a minimal, well commented example for writing a snap collector plugin in go. It advertises a single metric.
The example task collects this metric.

## Get the code:
`go get github.com/intelsdi-x/snap-plugin-collector-skeleton`

## Build:
```
cd $GOPATH/github.com/intelsdi-x/snap-plugin-collector-skeleton
go build
```

## Run
`Note: requires running snapd instance`

`/path/to/snapctl plugin load snap-plugin-collector-skeleton`


