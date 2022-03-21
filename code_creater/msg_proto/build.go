// Author: mike.qiao
// File:build
// Date:2022/3/16 10:35

package msg_proto

//go:generate protoc -I. -I../msg_proto --go_out=. --go_opt=paths=source_relative  msg.proto data.proto
