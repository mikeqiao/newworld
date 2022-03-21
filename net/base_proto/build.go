// Author: mike.qiao
// File:build
// Date:2022/3/18 11:20

package base_proto

//go:generate protoc -I. -I../base_proto --go_out=. --go_opt=paths=source_relative  base.proto
