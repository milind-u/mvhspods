module tests

go 1.20

replace mvhspods => ../mvhspods

require (
	github.com/milind-u/glog v0.0.0-20211106182138-9da3a6a0e251
	mvhspods v0.0.0-00010101000000-000000000000
)

require github.com/golang/glog v0.0.0-20210429001901-424d2337a529 // indirect
