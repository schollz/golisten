package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func matchParentheses(s string) string {
	i := strings.Index(s, "(")
	if i >= 0 {
		j := strings.Index(s[i:], ")")
		if j >= 0 {
			return "(" + s[i:i+j-1] + ")"
		}
	}
	return ""
}

func findFiles(searchDir string) (fileList []string, err error) {
	fileList = []string{}
	err = filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileName := filepath.Base(path)
		if !strings.Contains(fileName, ".mp3") {
			return nil
		}
		justFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		fileList = append(fileList, justFileName)
		return nil
	})
	return
}

func main() {

	router := gin.Default()
	router.StaticFS("/song", http.Dir("/home/zns/Music/music/"))
	router.StaticFS("/assets", http.Dir("assets"))
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		songFiles, err := findFiles("/home/zns/Music/music/")
		if err != nil {
			panic(err)
		}
		type SongList struct {
			ID   int
			Name string
			File string
		}
		songs := make([]SongList, len(songFiles))
		for i, song := range songFiles {
			songs[i].ID = i + 1
			songs[i].Name = song
			// Strip the name if its Youtube related
			if len(song) > 12 {
				if song[len(song)-12:len(song)-11] == "-" {
					songs[i].Name = song[:len(song)-12]
				}
			}
			// Strip parentheses from name
			hasP := matchParentheses(songs[i].Name)
			if len(hasP) > 0 {
				songs[i].Name = strings.Replace(songs[i].Name, hasP, "", -1)
			}
			songs[i].File = song
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
			"songs": songs,
		})
	})
	router.Run(":8080")
}
