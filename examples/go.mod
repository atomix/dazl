module github.com/atomix/dazl/examples

go 1.19

require (
	github.com/atomix/dazl v1.0.0
	github.com/atomix/dazl/zap v0.0.0-20230415091330-519194126903
)

require (
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/atomix/dazl => ../
replace github.com/atomix/dazl/zap => ../zap
