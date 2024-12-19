package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/wtolson/go-taglib"
)

// List of known audio file extensions
var audioExtensions = []string{
	".mp3", ".wav", ".flac", ".aac", ".ogg", ".wma", ".m4a", ".aiff", ".alac", ".opus",
}

type FileInfo struct {
	DirectoryPath    string
	FileName         string
	FileExtension    string
	Title            string
	Artist           string
	Album            string
	Genre            string
	Year             int
	Bitrate          int
	Samplerate       int
	Channels         int
	Length           time.Duration
	Track            int
	CreationDate     time.Time
	ModificationDate time.Time
}

// Checks if the file has a known audio extension
func isAudioFile(extension string) bool {
	extension = strings.ToLower(extension)
	for _, ext := range audioExtensions {
		if ext == extension {
			return true
		}
	}
	return false
}

// getCreationDate retrieves the creation date (birth time) of a file using syscall.Stat_t.
func getCreationDate(path string) (time.Time, error) {
	var stat syscall.Stat_t
	if err := syscall.Stat(path, &stat); err != nil {
		return time.Time{}, err
	}
	// Use Ctimespec as a reliable fallback for file creation time
	return time.Unix(stat.Ctimespec.Sec, stat.Ctimespec.Nsec), nil
}

func scanDirectory(root string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fileExt := strings.ToLower(filepath.Ext(info.Name()))
			if isAudioFile(fileExt) {
				dirPath := filepath.Dir(path)
				fileName := info.Name()
				modDate := info.ModTime()

				// Get creation date
				creationDate, err := getCreationDate(path)
				if err != nil {
					creationDate = time.Time{} // Use zero time if unavailable
				}

				// Open the audio file
				fullpath := filepath.Join(dirPath, fileName)
				audioMetadata, err := taglib.Read(fullpath)
				if err != nil {
					log.Fatal(err)
				}
				// Extract metadata
				title := audioMetadata.Title()
				artist := audioMetadata.Artist()
				album := audioMetadata.Album()
				year := audioMetadata.Year()
				genre := audioMetadata.Genre()
				bitrate := audioMetadata.Bitrate()
				samplerate := audioMetadata.Samplerate()
				channels := audioMetadata.Channels()
				length := audioMetadata.Length()
				track := audioMetadata.Track()

				files = append(files, FileInfo{
					DirectoryPath:    dirPath,
					FileName:         fileName,
					FileExtension:    fileExt,
					CreationDate:     creationDate,
					ModificationDate: modDate,
					Title:            title,
					Artist:           artist,
					Album:            album,
					Year:             year,
					Genre:            genre,
					Bitrate:          bitrate,
					Samplerate:       samplerate,
					Channels:         channels,
					Length:           length,
					Track:            track,
				})
			}
		}
		return nil
	})

	return files, err
}

func main() {
	// Get directory from command-line arguments or default to current directory
	root := "./"
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	// Scan the directory
	files, err := scanDirectory(root)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Display results
	for _, file := range files {
		fmt.Printf("Directory: %s\n", file.DirectoryPath)
		fmt.Printf("File Name: %s\n", file.FileName)
		fmt.Printf("File Extension: %s\n", file.FileExtension)
		fmt.Printf("Creation Date: %s\n", file.CreationDate)
		fmt.Printf("Modification Date: %s\n", file.ModificationDate)

		fmt.Printf("Title: %s\n", file.Title)
		fmt.Printf("Artist: %s\n", file.Artist)
		fmt.Printf("Album: %s\n", file.Album)
		fmt.Printf("Year: %d\n", file.Year)
		fmt.Printf("Genre: %s\n", file.Genre)
		fmt.Println("-------------")
	}
}
