// Author: mike.qiao
// File:data
// Date:2022/3/14 16:05

package data

type MemoryData interface {
	Init()
	Update()
	Close()
}
