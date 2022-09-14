package filesystem

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"time"

	"github.com/avnes/test-persistence/pkg/common"
)

const appFilesDir = "/tmp/files"

type metadata struct {
	Name     string    `json:"filename"`
	Size     int64     `json:"size_in_bytes"`
	Elapsed  float64   `json:"elapsed_in_seconds"`
	Modified time.Time `json:"last_modified"`
	Type     string    `json:"type"`
}

func (f *metadata) record(filename string, size int64, elapsed float64, modified time.Time) metadata {
	f.Name = filename
	f.Size = size
	f.Elapsed = elapsed
	f.Modified = modified
	f.Type = filename[len(filename)-4:]
	return *f
}

func getDirectory() string {
	directory, exist := os.LookupEnv("APP_FILES_DIR")
	if !exist {
		directory = appFilesDir
	}
	return directory
}

func getFilesInDirectory() ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(getDirectory())
	if err != nil {
		return nil, err
	}
	return files, nil
}

func getElapsedTime(filename string) (float64, error) {
	metaFile := getDirectory() + "/" + filename[0:len(filename)-4] + "meta"
	jsonFile, err := os.Open(metaFile)
	if err != nil {
		return 0.0, err
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return 0.0, err
	}
	var fs metadata
	json.Unmarshal(byteValue, &fs)
	return fs.Elapsed, err
}

func httpGetFiles() ([]metadata, error) {
	files, err := getFilesInDirectory()
	if err != nil {
		return nil, err
	}

	var jsonFragment []metadata

	for _, file := range files {
		// We only need to show the data files. Not the meta files.
		if file.Name()[len(file.Name())-4:] == "data" {
			elapsed, err := getElapsedTime(file.Name())
			if err != nil {
				return nil, err
			}
			record := new(metadata).record(file.Name(), file.Size(), elapsed, file.ModTime())
			jsonFragment = append(jsonFragment, record)
		}
	}
	return jsonFragment, nil
}

func httpGetFilesCount() ([]common.Counter, error) {
	files, err := getFilesInDirectory()
	if err != nil {
		return nil, err
	}
	var jsonFragment []common.Counter
	record := new(common.Counter).Record(len(files) / 2) // Divide by 2 because we don't care about *.meta files
	jsonFragment = append(jsonFragment, record)
	return jsonFragment, nil
}

func writeBinaryFile(filename string, size int64) (bool, error) {
	data := make([]byte, int64(size), int64(size))
	f, err := os.Create(fmt.Sprintf("%s/%s", getDirectory(), filename))
	if err != nil {
		return false, err
	}
	defer f.Close()
	bytesWritten, err := f.Write(data)
	_ = bytesWritten
	if err != nil {
		return false, err
	}
	return true, nil
}

func writeFileStats(filename string, size int64, elapsed time.Duration) {
	content := metadata{
		Name:     filename,
		Size:     size,
		Elapsed:  elapsed.Seconds(),
		Modified: time.Now(),
	}
	jsonContent, _ := json.MarshalIndent(content, "", " ")
	metaFile := getDirectory() + "/" + filename[0:len(filename)-5] + ".meta"
	_ = ioutil.WriteFile(metaFile, jsonContent, 0644) // Disregard errors on writing metadata
}

func httpPostFiles(count int, size int64) ([]metadata, error) {
	var jsonFragment []metadata
	for i := 0; i < count; i++ {
		start := time.Now()
		filename := fmt.Sprintf("%s.data", string(start.Format(time.RFC3339Nano)))
		_, err := writeBinaryFile(filename, size)
		if err != nil {
			return jsonFragment, err
		}
		elapsed := time.Since(start)
		writeFileStats(filename, size, elapsed)
		record := new(metadata).record(filename, size, elapsed.Seconds(), start)
		jsonFragment = append(jsonFragment, record)
	}
	return jsonFragment, nil
}

func CreateDirectory() error {
	_, err := os.Stat(getDirectory())
	if os.IsNotExist(err) {
		err := os.MkdirAll(getDirectory(), 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func failIfDirectoryDoesNotExist() error {
	_, err := os.Stat(getDirectory())
	if os.IsNotExist(err) {
		return err
	}
	return nil
}

func httpDeleteFiles() (bool, error) {
	files, err := getFilesInDirectory()
	if err != nil {
		return false, err
	}
	if len(files) == 0 {
		return false, nil
	}
	for _, file := range files {
		err := os.Remove(getDirectory() + "/" + file.Name())
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
