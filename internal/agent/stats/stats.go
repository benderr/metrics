// Package stats contain DTO for collect stats
package stats

type Item struct {
	Name  string
	Type  string
	Delta int64
	Value float64
}
