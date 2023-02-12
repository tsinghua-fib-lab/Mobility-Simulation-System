package parallel

import (
	"sync"

	"github.com/samber/lo"
)

// 自行wg.Wait()
func GoForWithWaitGroup[T any](values []T, chunkSize int, wg *sync.WaitGroup, f func(v T)) {
	if chunkSize < 1 {
		panic("chunk size >= 1")
	}
	if chunkSize == 1 {
		wg.Add(len(values))
		for _, vv := range values {
			v := vv
			go func() {
				defer wg.Done()
				f(v)
			}()
		}
	} else {
		chunks := lo.Chunk(values, chunkSize)
		wg.Add(len(chunks))
		for _, chunk := range chunks {
			c := chunk
			go func() {
				defer wg.Done()
				for _, v := range c {
					f(v)
				}
			}()
		}
	}
}

func GoFor[T any](values []T, chunkSize int, f func(v T)) {
	var wg sync.WaitGroup
	GoForWithWaitGroup(values, chunkSize, &wg, f)
	wg.Wait()
}
