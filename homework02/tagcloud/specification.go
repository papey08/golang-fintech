package tagcloud

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	// tagStats is a slice of TagStat structs ordered by OccurenceCount descending
	tagStats []TagStat

	// indexes is a map where key is a tag and value is its index in tagStats
	indexes map[string]int
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
func New() TagCloud {
	return TagCloud{
		tagStats: make([]TagStat, 0),
		indexes:  make(map[string]int),
	}
}

// AddTag adds a new tag in TagCloud if its not in it or increasing
// OccurenceCount of already existing tag and moving it on its new place if needed.
// Time complexity: O(len(tagStats)) in worst case
func (t *TagCloud) AddTag(tag string) {
	if ind, ok := t.indexes[tag]; !ok { // check if tag is already in TagCloud
		t.indexes[tag] = len(t.tagStats)
		t.tagStats = append(t.tagStats, TagStat{tag, 1})
	} else {
		t.tagStats[ind].OccurrenceCount++
		// keeping tagStats odered by switching current tag with his left neighbour
		for ind != 0 && t.tagStats[ind-1].OccurrenceCount <= t.tagStats[ind].OccurrenceCount {
			tag1 := t.tagStats[ind]
			tag2 := t.tagStats[ind-1]
			t.indexes[tag1.Tag], t.indexes[tag2.Tag] = t.indexes[tag2.Tag], t.indexes[tag1.Tag]
			t.tagStats[ind-1], t.tagStats[ind] = t.tagStats[ind], t.tagStats[ind-1]
			ind--
		}
	}
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
func (t *TagCloud) TopN(n int) []TagStat {
	if n <= len(t.tagStats) {
		return t.tagStats[:n]
	} else {
		return t.tagStats
	}
}
