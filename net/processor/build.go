// Author: mike.qiao
// File:build
// Date:2022/3/10 16:46

package processor

//go:generate protoc -I. -I../base_proto --go_out=. --go_opt=paths=source_relative  base.proto
