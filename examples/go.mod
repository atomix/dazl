module github.com/atomix/dazl/examples

go 1.19

require (
	github.com/atomix/dazl v1.1.0
	github.com/atomix/dazl/zerolog v0.0.0-20230415093331-60e3c8575eb5
)

require (
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/rs/zerolog v1.29.0 // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/atomix/dazl => ../

replace github.com/atomix/dazl/zap => ../zap

replace github.com/atomix/dazl/zerolog => ../zerolog
