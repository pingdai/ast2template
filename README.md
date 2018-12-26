# ast2template
A Lightweight tools By Golang AST.Now it has supported below:
- [gorm model]()

## Usage
This is an executable program,if you want to use it,you need to follow steps:
1. Download it:

```sh
go get -u github.com/pingdai/ast2template
```

or

```sh
git clone https://github.com/pingdai/ast2template.git
```

2. Install it:

```sh
go install
```

3. You maybe get some other packages:

```sh
go get -u github.com/spf13/cobra
go get -u github.com/jinzhu/gorm
```

## Quick start

```sh
$ ast2template
```

```sh
ast tools

Usage:
  ast2template [command]

Available Commands:
  gen         generators
  help        Help about any command

Flags:
  -h, --help   help for ast2template

Use "ast2template [command] --help" for more information about a command.

```

See detail useage,click [here](https://github.com/pingdai/ast2template/blob/master/examples/student.go)
## The third party relying on package
- [cobra](https://github.com/spf13/cobra)
- [gorm](https://github.com/jinzhu/gorm)