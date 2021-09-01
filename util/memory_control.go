package util

import memory_retention "github.com/yuto51942/memory-retention"

// Clear all memory.
func ClearAll() {
	memory_retention.DeleteAll()
}
