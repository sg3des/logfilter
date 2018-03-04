# logfilter

logfilter is inspired by [logutils](https://github.com/hashicorp/logutils), has the same principle.

LogLevel replaced by LogPrefix, logfilter does not have any levels, its operate prefixes. logfilter allow add multiple writers with different prefixes. 

### install

```go
go get -u github.com/sg3des/logilter
```

## usage

```go
filter := logfilter.NewFilter()
filter.AddWriter(os.Stdout, "DEBUG", "INFO", "ERROR")
filter.AddStrictWriter(file, "ERROR", "FILE")

log.SetOutput(filter)
log.Print("information")
log.Print("[DEBUG] debug information")
log.Print("[INFO] some information")
log.Print("[ERROR] error writed to stdout and to file")
log.Print("[FILE] info writed to stdout and to file")
```