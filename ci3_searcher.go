package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Vector100 [100]float32

type Item struct {
	Type     string
	Name     string
	Path     string
	Vector   Vector100
	Metadata string
}

func (v Vector100) Dot(other Vector100) float32 {
	var sum float32
	for i := 0; i < 100; i++ {
		sum += v[i] * other[i]
	}
	return sum
}

func (v *Vector100) Normalize() {
	var sum float32
	for _, val := range v {
		sum += val * val
	}
	norm := float32(math.Sqrt(float64(sum)))
	if norm > 0 {
		for i := 0; i < 100; i++ {
			v[i] /= norm
		}
	}
}

func Vectorize(text string, itemType string, name string) Vector100 {
	var v Vector100
	
	// Add type and name to content for better weight
	content := fmt.Sprintf("%s %s %s", itemType, name, text)
	tokens := strings.Fields(strings.ToLower(content))
	
	for _, token := range tokens {
		h := fnv.New32a()
		h.Write([]byte(token))
		idx := h.Sum32() % 100
		v[idx] += 1.0
	}

	v.Normalize()
	return v
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ci3_searcher.go <path_to_ci3_project>")
		return
	}

	root := os.Args[1]
	absRoot, _ := filepath.Abs(root)
	fmt.Printf("Indexing CI3 project at: %s... 0 items indexed", absRoot)
	var indexedCount int
	var index []Item

	err := filepath.Walk(absRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Propagate errors
		}
		if info.IsDir() {
			name := info.Name()
			// Skip specific system/dependency directories
			if name == ".git" || name == "logs" || name == "node_modules" || name == "system" || name == "vendor" {
				return filepath.SkipDir // Skip these directories entirely
			}
			return nil // Continue walking into other directories
		}
		// If it's a file, process it
		indexedCount++
		if indexedCount%100 == 0 { // Update every 100 files
			fmt.Printf("\rIndexing CI3 project at: %s... %d items indexed", absRoot, indexedCount)
		}

		name := info.Name()
		rel, _ := filepath.Rel(absRoot, path)
		relLower := strings.ToLower(rel)

		if strings.Contains(relLower, ".git") || strings.Contains(relLower, "logs") {
			return nil
		}
		
		itemType := "other"
		if strings.Contains(relLower, "controllers/") {
			itemType = "controller"
		} else if strings.Contains(relLower, "models/") {
			itemType = "model"
		} else if strings.Contains(relLower, "views/") {
			itemType = "view"
		} else if strings.Contains(relLower, "config/") {
			itemType = "config"
		}

		content, _ := os.ReadFile(path)
		contentText := string(content)
		index = append(index, Item{
			Type:     itemType,
			Name:     name,
			Path:     rel,
			Metadata: contentText,
			Vector:   Vectorize(contentText, itemType, name),
		})
		return nil
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Successfully indexed %d items.\n", len(index))
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\nSearch> ")
		if !scanner.Scan() {
			break
		}
		query := scanner.Text()
		if query == "q" {
			break
		}
		
		qVec := Vectorize(query, "", "")
		type Result struct {
			Score float32
			Item  Item
		}
		var results []Result

		for _, item := range index {
			score := qVec.Dot(item.Vector)
			if score > 0.01 { // Lowered threshold
				results = append(results, Result{Score: score, Item: item})
			}
		}

		sort.Slice(results, func(i, j int) bool {
			return results[i].Score > results[j].Score
		})

		if len(results) == 0 {
			fmt.Println("No matches found.")
			continue
		}

		for i := 0; i < len(results) && i < 10; i++ {
			r := results[i]
			fmt.Printf("%.4f | %-12s | %s\n", r.Score, r.Item.Type, r.Item.Path)
		}
	}
}
