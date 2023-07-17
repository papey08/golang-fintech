package tagcloud_test

import (
	"fmt"
	"lecture02_homework/tagcloud"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyTagCloud(t *testing.T) {
	tc := tagcloud.New()
	topN := tc.TopN(1000)
	assert.Len(t, topN, 0, "empty tag cloud returned non-empty topN elements %v", topN)
}

func TestTopNGreaterThanCloudSize(t *testing.T) {
	tc := tagcloud.New()
	tc.AddTag("t1")

	requestCount := 10
	topN := tc.TopN(requestCount)
	assert.Len(t, topN, 1, "TopN(%d) returned array with invalid size: %v", requestCount, topN)
}

func TestHappyPath(t *testing.T) {
	tc := tagcloud.New()

	tc.AddTag("single-occurrence")
	tc.AddTag("multiple-occurrence")
	tc.AddTag("multiple-occurrence")

	top := tc.TopN(1)
	assert.Len(t, top, 1, "TopN(1) returned %d elements", len(top))

	if assert.Equal(t, "multiple-occurrence", top[0].Tag) {
		assert.Equal(t, 2, top[0].OccurrenceCount)
	}
}

func TestTopN(t *testing.T) {
	tc := tagcloud.New()
	size := 1000
	for i := 0; i < size; i++ {
		for j := 0; j < i; j++ {
			tc.AddTag(fmt.Sprintf("%d", i))
		}
	}

	validateTopN := func(n int) {
		topN := tc.TopN(n)
		assert.Len(t, topN, n)

		for i, el := range topN {
			value := size - i - 1
			tagName := fmt.Sprintf("%d", value)
			assert.Equal(t, tagName, el.Tag, "TopN(%d) returned elements in wrong order (bad tag name at %d): %v", n, i, topN)
			assert.Equal(t, value, el.OccurrenceCount, "TopN(%d) returned elements in wrong order (bad occurrence count at %d): %v", n, i, topN)
		}
	}

	for i := 0; i < size; i++ {
		validateTopN(i)
	}
}

func TestTopNWithRepeatedOccurrence(t *testing.T) {
	tc := tagcloud.New()
	tc.AddTag("t1")
	tc.AddTag("t2")
	tc.AddTag("t3")
	tc.AddTag("t4")

	requestCount := 3
	topN := tc.TopN(requestCount)
	assert.Len(t, topN, requestCount)

	distinctMap := make(map[string]struct{})
	for _, v := range topN {
		assert.Equal(t, 1, v.OccurrenceCount)
		distinctMap[v.Tag] = struct{}{}
	}

	assert.Len(t, distinctMap, requestCount, "TopN(%d) returned array with non-distinct tags: %v", requestCount, topN)
}

// TestTopNDifficult checks if tags replaces correct if they be adding in random order
func TestTopNWithRandomAddingOrder(t *testing.T) {
	tags := []tagcloud.TagStat{
		{
			Tag:             "tag01",
			OccurrenceCount: 50,
		},
		{
			Tag:             "tag02",
			OccurrenceCount: 50,
		},
		{
			Tag:             "tag03",
			OccurrenceCount: 100,
		},
		{
			Tag:             "tag04",
			OccurrenceCount: 100,
		},
		{
			Tag:             "tag05",
			OccurrenceCount: 100,
		},
		{
			Tag:             "tag06",
			OccurrenceCount: 100,
		},
		{
			Tag:             "tag07",
			OccurrenceCount: 100,
		},
		{
			Tag:             "tag08",
			OccurrenceCount: 200,
		},
		{
			Tag:             "tag09",
			OccurrenceCount: 200,
		},
		{
			Tag:             "tag10",
			OccurrenceCount: 200,
		},
	}

	tc := tagcloud.New()

	wg := new(sync.WaitGroup)
	mu := new(sync.Mutex)
	for _, tag := range tags {
		wg.Add(1)
		go func(tag tagcloud.TagStat) {
			defer wg.Done()
			for i := 0; i < tag.OccurrenceCount-1; i++ {
				mu.Lock()
				tc.AddTag(tag.Tag)
				mu.Unlock()
			}
		}(tag)
	}
	wg.Wait()
	for _, tag := range tags {
		tc.AddTag(tag.Tag)
	}

	n := 5
	correctTopN := []tagcloud.TagStat{
		{
			Tag:             "tag10",
			OccurrenceCount: 200,
		},
		{
			Tag:             "tag09",
			OccurrenceCount: 200,
		},
		{
			Tag:             "tag08",
			OccurrenceCount: 200,
		},
		{
			Tag:             "tag07",
			OccurrenceCount: 100,
		},
		{
			Tag:             "tag06",
			OccurrenceCount: 100,
		},
	}
	gotTopN := tc.TopN(n)

	assert.Len(t, gotTopN, n, "TopN(%d) returned array: %v", n, gotTopN)
	for i := range gotTopN {
		assert.Equal(t, correctTopN[i].Tag, gotTopN[i].Tag, "TopN(%d)[%d] wrong tag: %s", n, i, gotTopN[i].Tag)
		assert.Equal(t, correctTopN[i].OccurrenceCount, gotTopN[i].OccurrenceCount, "TopN(%d)[%d] wrong OccurrenceCount: %s", n, i, gotTopN[i].OccurrenceCount)
	}
}
