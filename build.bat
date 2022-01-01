echo $(date)
go build -tags=jsoniter -ldflags="-X 'main.BuildTime=%date% - %time:~0,8%'" .