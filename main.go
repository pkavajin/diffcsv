package main

import (
	"bufio"
	"flag"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	var aPath string
	var bPath string
	var newPath string
	var samePath string
	var deletedPath string
	var hasHeader bool
	flag.StringVar(&aPath, "a", "", "first file")
	flag.StringVar(&bPath, "b", "", "second file")
	flag.StringVar(&newPath, "out-added", "added.csv", "file where the rows should be written which are in b but not in a (aka: added)")
	flag.StringVar(&samePath, "out-both", "both.csv", "file where the rows should be written which are in both b and a (aka: equal")
	flag.StringVar(&deletedPath, "out-deleted", "deleted.csv", "file where the rows should be written which are in a but not in b (aka: deleted")
	flag.BoolVar(&hasHeader, "header", true, "file contains header")
	flag.Parse()

	logrus.SetLevel(logrus.DebugLevel)

	start := time.Now()
	aHeader, aHashlist, aErr := readCSV(aPath, hasHeader)
	if aErr != nil {
		logrus.WithFields(logrus.Fields{
			"file":   aPath,
			"reason": aErr.Error(),
		}).Fatal("couldn't read a csv")
	}
	logrus.WithFields(logrus.Fields{
		"lat":  time.Since(start).String(),
		"file": aPath,
	}).Info("read file")

	start = time.Now()
	bHeader, bHashlist, bErr := readCSV(bPath, hasHeader)
	if aErr != nil {
		logrus.WithFields(logrus.Fields{
			"file":   bPath,
			"reason": bErr.Error(),
		}).Fatal("couldn't read b csv")
	}
	logrus.WithFields(logrus.Fields{
		"lat":  time.Since(start).String(),
		"file": bPath,
	}).Info("read file")

	if aHeader != bHeader {
		logrus.WithFields(logrus.Fields{
			"aHeader": aHeader,
			"bHeader": bHeader,
		}).Fatal("header definition changed - everything is different")
	}

	var newFile *os.File
	var sameFile *os.File
	var deletedFile *os.File
	var newWriter *bufio.Writer
	var sameWriter *bufio.Writer
	var deletedWriter *bufio.Writer
	var err error

	if newPath != "" {
		newFile, err = os.Create(newPath)
		if err != nil {
			logrus.WithField("file", newPath).Fatal(err)
		}

		if hasHeader {
			newFile.WriteString(aHeader + "\n")
		}

		newWriter = bufio.NewWriter(newFile)
		defer newWriter.Flush()
		defer newFile.Sync()
		defer newFile.Close()
	}

	if samePath != "" {
		sameFile, err = os.Create(samePath)
		if err != nil {
			logrus.WithField("file", samePath).Fatal(err)
		}

		if hasHeader {
			sameFile.WriteString(aHeader + "\n")
		}

		sameWriter = bufio.NewWriter(sameFile)
		defer sameWriter.Flush()
		defer sameFile.Sync()
		defer sameFile.Close()
	}

	if deletedPath != "" {
		deletedFile, err = os.Create(deletedPath)
		if err != nil {
			logrus.WithField("file", deletedPath).Fatal(err)
		}

		if hasHeader {
			deletedFile.WriteString(aHeader + "\n")
		}

		deletedWriter = bufio.NewWriter(deletedFile)
		defer deletedWriter.Flush()
		defer deletedFile.Sync()
		defer deletedFile.Close()
	}

	start = time.Now()

	additions := 0
	same := 0
	deletions := 0

	for line := range bHashlist {
		_, isInA := aHashlist[line]
		if isInA {
			same++
			if samePath != "" {
				sameWriter.WriteString(line + "\n")
			}
		} else {
			additions++
			if newPath != "" {
				newWriter.WriteString(line + "\n")
			}
		}
	}

	for line := range aHashlist {
		_, isInB := bHashlist[line]
		if !isInB {
			deletions++
			if deletedPath != "" {
				deletedWriter.WriteString(line + "\n")
			}
		}
	}

	if newPath != "" {
		newWriter.Flush()
		newFile.Sync()
		newFile.Close()
	}

	if samePath != "" {
		sameWriter.Flush()
		sameFile.Sync()
		sameFile.Close()
	}

	if deletedPath != "" {
		deletedWriter.Flush()
		deletedFile.Sync()
		deletedFile.Close()
	}

	logrus.WithFields(logrus.Fields{
		"lat":       time.Since(start).String(),
		"additions": additions,
		"same":      same,
		"deletions": deletions,
		"totalA":    len(aHashlist),
		"totalB":    len(bHashlist),
	}).Info("Finished diff")
}

func readCSV(path string, header bool) (string, map[string]struct{}, error) {
	inputFile, err := os.Open(path)
	if err != nil {
		return "", nil, err
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	headerStr := ""
	i := 0

	hashlist := make(map[string]struct{})

	for scanner.Scan() {
		line := scanner.Text()
		if i == 0 && header {
			headerStr = line
			i++
			continue
		}

		hashlist[line] = struct{}{}

		i++
	}

	if err := scanner.Err(); err != nil {
		return "", nil, err
	}

	return headerStr, hashlist, nil
}
