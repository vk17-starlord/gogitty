package core

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"gogitty/pkg/utils"
	"io"
	"os"
	"strconv"
	"strings"
)

// GitObject interface
type GitObject interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) error
	Init() error
	Format() string
}

func ObjectRead(repo string, sha string) {
	path := utils.RepoPath(repo, "objects", sha[:2], sha[2:])

	objectFile, err := os.ReadFile(path)
	utils.Check(err)

	// use zlib for reading content of the objectFILE
	data, err := zlib.NewReader(bytes.NewReader(objectFile))
	if err != nil {
		utils.Check(err)
	}

	rawContent, err := io.ReadAll(data)
	utils.Check(err)

	// Find the object type in the rawContent
	index := bytes.Index(rawContent, []byte(" "))

	// extract the type and size of the object
	y := bytes.Index(rawContent, []byte("\x00"))
	contentSize := string(rawContent[index:y])
	contentSize = strings.TrimSpace(contentSize)
	size, _ := strconv.Atoi(contentSize)

	// extract the rawContent of the object
	Actualcontent := rawContent[y+1:]
	if len(Actualcontent) == size {

		fmt.Print(string(Actualcontent))

	} else {

		fmt.Println("Content is not correct")
	}

	// close the reader
	defer data.Close()

}

func ObjectWrite(gitobject GitObject, repo string, write ...bool) (string, error) {

	// set default value of write to true
	writeFlag := true
	if len(write) > 0 {
		writeFlag = write[0]
	}

	// serialize the object
	data, _ := gitobject.Serialize()

	// read the file
	result := gitobject.Format() + " " + strconv.Itoa(len(data)) + "\x00" + string(data)

	hash := sha1.New()
	hash.Write([]byte(result))
	sha := fmt.Sprintf("%x", hash.Sum(nil))

	if writeFlag {
		Objectpath := utils.RepoPath(repo, "objects", sha[0:2])
		utils.EnsureDir(Objectpath)
		var buffer bytes.Buffer
		w := zlib.NewWriter(&buffer)
		w.Write([]byte(result))
		w.Close()
		err := os.WriteFile(Objectpath+"/"+sha[2:], buffer.Bytes(), 0644)

		if err != nil {
			return "", err
		}

	}

	return sha, nil

}
