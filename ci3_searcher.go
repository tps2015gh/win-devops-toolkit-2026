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

type Vector300 [300]float32

type Item struct {
	Type     string
	Name     string
	Path     string
	Vector   Vector300
	Metadata string
}

func (v Vector300) Dot(other Vector300) float32 {
	var sum float32
	for i := 0; i < 300; i++ {
		sum += v[i] * other[i]
	}
	return sum
}

func (v *Vector300) Normalize() {
	var sum float32
	for _, val := range v {
		sum += val * val
	}
	norm := float32(math.Sqrt(float64(sum)))
	if norm > 0 {
		for i := 0; i < 300; i++ {
			v[i] /= norm
		}
	}
}

func Vectorize(text string, itemType string, name string) Vector300 {
	var v Vector300
	
	// Add type and name to content for better weight
	content := fmt.Sprintf("%s %s %s", itemType, name, text)
	tokens := strings.Fields(strings.ToLower(content))
	
	for _, token := range tokens {
		h := fnv.New32a()
		h.Write([]byte(token))
		idx := h.Sum32() % 300
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

	// Display current vector dimension
	fmt.Printf("CI3 Searcher initialized with a %d-dimensional vector space.\n", 300)
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

		ext := strings.ToLower(filepath.Ext(path))
		rel, _ := filepath.Rel(absRoot, path)
		relLower := strings.ToLower(rel)
		
		// Skip binary and image files
		skipExts := map[string]bool{
			".pdf": true, ".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true, ".tiff": true, ".webp": true,
			".exe": true, ".zip": true, ".rar": true, ".7z": true, ".tar": true, ".gz": true, ".bz2": true,
			".iso": true, ".dmg": true, ".bin": true, ".dat": true, ".db": true, ".sqlite": true,
			".woff": true, ".woff2": true, ".ttf": true, ".otf": true, ".eot": true, // Fonts
			".mp3": true, ".wav": true, ".ogg": true, ".flac": true, ".aac": true, // Audio
			".mp4": true, ".avi": true, ".mkv": true, ".mov": true, ".wmv": true, // Video
		}
		if skipExts[ext] {
			return nil
		}

		// Also skip git, logs, and other non-app files that were not caught by directory skip
		if strings.Contains(relLower, ".git") || strings.Contains(relLower, "logs") || strings.Contains(relLower, "node_modules") || strings.Contains(relLower, "cache") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil // Silently skip unreadable files
		}
		contentText := string(content)
		
		itemType := "other"
		
		// CI3 Component Detection
		if strings.Contains(relLower, "controllers/") && ext == ".php" {
			itemType = "controller"
		} else if strings.Contains(relLower, "models/") && ext == ".php" {
			itemType = "model"
		} else if strings.Contains(relLower, "views/") && (ext == ".php" || ext == ".html" || ext == ".phtml") {
			itemType = "view"
		} else if strings.Contains(relLower, "config/") && ext == ".php" {
			itemType = "config"
		} else if ext == ".js" {
			itemType = "js"
		} else if ext == ".sql" { // Explicitly add SQL if needed, but not common CI3 app files
			itemType = "sql"
		}

		// Only index files of recognized types or generic PHP/HTML/JS/SQL
		if itemType == "other" && ext != ".php" && ext != ".html" && ext != ".phtml" && ext != ".js" && ext != ".sql" {
			return nil // Skip this file if it's not a recognized CI3 component and not a general code/script file
		}

		// Simplified database interaction detection for models/controllers
		if (itemType == "model" || itemType == "controller") && strings.Contains(contentText, "$this->db->") {
			itemType += "/db"
		}

		// Correctly get the base name from info
		baseName := info.Name()

		index = append(index, Item{
			Type:     itemType,
			Name:     baseName,
			Path:     rel,
			Metadata: contentText,
			Vector:   Vectorize(contentText, itemType, baseName),
		})
		return nil
	}) // Correctly close filepath.Walk anonymous function

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

		// Display up to 20 results, export to file if more
		displayLimit := 20
		if len(results) > displayLimit {
			outFile := "search_results.txt"
			f, err := os.Create(outFile)
			if err == nil {
				fmt.Fprintf(f, "Search Query: %s\nTotal Results: %d\n\n", query, len(results))
				for _, r := range results {
					fmt.Fprintf(f, "[Score: %.4f] [%s] %s\nPath: %s\n\n", r.Score, r.Item.Type, r.Item.Name, r.Item.Path)
				}
				f.Close()
				fmt.Printf("Found %d results. Showing top %d. Full list exported to %s\n", len(results), displayLimit, outFile)
			} else {
				fmt.Printf("Found %d results. Showing top %d. Could not write to file: %v\n", len(results), displayLimit, err)
			}
		} else {
			fmt.Printf("Found %d results:\n", len(results))
		}


		for i := 0; i < len(results) && i < displayLimit; i++ {
			r := results[i]
			fmt.Printf("%.4f | %-12s | %s\n", r.Score, r.Item.Type, r.Item.Path)
		}
	}
}
