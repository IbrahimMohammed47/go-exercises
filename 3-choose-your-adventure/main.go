package main

// make http server
import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Chapter struct {
	Title   string       `json:"title"`
	Story   []string     `json:"story"`
	Options []ChapterArc `json:"options"`
}
type ChapterArc struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func main() {

	storyHandler := getStoryHandler("gopher.json")
	http.HandleFunc("/story", storyHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

// "gopher.json"
func getStoryHandler(storiesFile string) http.HandlerFunc {
	jsonFile, err := os.Open(storiesFile)
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var jsonMap map[string]Chapter
	json.Unmarshal([]byte(byteValue), &jsonMap)

	jsonFile.Close()
	tmplt, _ := template.ParseFiles("story.html")

	storyHandler := func(w http.ResponseWriter, r *http.Request) {
		chapterKey := r.URL.Query().Get("chapter")
		if chapterKey == "" {
			chapterKey = "intro"
		}
		chapter, ok := jsonMap[chapterKey]
		if !ok {
			http.NotFound(w, r)
			return
		}

		err := tmplt.Execute(w, chapter)
		if err != nil {
			return
		}
	}
	return storyHandler
}
