# Executable to import Raster to CartoDB dataset

## Install

Install go: [Install](https://golang.org/doc/install)

Install dependencies:

````
go get github.com/jessevdk/go-flags
go get github.com/Sirupsen/logrus
go get github.com/agonzalezro/cartodb_go
go get gopkg.in/cheggaaa/pb.v1
```

## Execute

In the project path, execute:

```
go run main.go <options>
```

## Generate executable

In the project path, execute:

````
go build . 
```