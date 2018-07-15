package main

import (
	"archive/tar"
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// TODO
type formatType string

const (
	formatPosix formatType = "posix"
	formatGnu   formatType = "gnu"
)

var (
	errFileNotExistTpl = "No Such File or Directory: %s"
	errFileAlreadyExistTpl = "File or Directory Already Exist: %s"
)

var (
	stdin       bool   = false   // get file list from stdin
	delim       string = "\n"    // delimiter
	filePath    string = ""      // path
	format      string = "posix" // TODO
	absPath     bool   = false
	compress    string = ""    // TODO
	ignoreError bool   = false // TODO
	outputFile  string = ""
)

type GlobalFlag struct {
	stdin    bool
	delim    byte
	format   string
	compress string
	absPath  bool
}

func init() {
	flag.BoolVar(&stdin, "stdin", false, "Get file list from stdin")
	flag.StringVar(&delim, "delim", "\n", "File list delimiter")
	flag.StringVar(&filePath, "file-path", "", "File path")
	flag.StringVar(&format, "format", "posix", "Create archive of the given format")
	flag.StringVar(&compress, "compress", "", "Compress type")
	flag.BoolVar(&absPath, "abs-path", false, "Archive file with abs path")
	flag.BoolVar(&ignoreError, "ignore-error", false, "Continue while error ocurre")
	flag.StringVar(&outputFile, "output-file", "", "Output File")
}

func archiveTar(fpList []string, outputFname string) (err error) {
	var fw *os.File
	fw, err = os.Create(outputFname)
	if err != nil {
		return
	}
	tw := tar.NewWriter(fw)

	for _, fp := range fpList {
		var f *os.File
		var fi os.FileInfo
		f, err = os.Open(fp)
		if err != nil {
			return
		}
		defer f.Close()
		fi, err = f.Stat()
		fstat := fi.Sys().(*syscall.Stat_t)
		err = tw.WriteHeader(&tar.Header{
			Name:       strings.TrimLeft(fp, string(filepath.Separator)),
			Size:       fi.Size(),
			ModTime:    fi.ModTime(),
			Mode:       int64(fi.Mode()),
			Uid:        int(fstat.Uid),
			Gid:        int(fstat.Gid),
			AccessTime: time.Unix(fstat.Atim.Sec, fstat.Atim.Nsec),
			ChangeTime: time.Unix(fstat.Ctim.Sec, fstat.Ctim.Nsec),
		})
		if err != nil {
			panic(err)
		}
		var writed int64
		if writed, err = io.Copy(tw, f); err != nil {
			if err == io.EOF && writed == fi.Size() {
				if err = tw.Flush(); err != nil {
					panic(err)
				}
				continue
			}
			panic(err)
		}
	}
	tw.Close()
	return
}

func verifyFilePath(fp string, abs bool) (fn string, err error) {
	fp = strings.TrimRight(fp, "\n")
	if fp == "" {
		err = fmt.Errorf(errFileNotExistTpl, fp)
		return
	}
	if abs {
		return filepath.Abs(filepath.Clean(fp))
	}
	return filepath.Clean(fp), err
}

func main() {
	flag.Parse()
	var err error

	// verify outputFile
	_, err = os.Stat(outputFile)
	if !os.IsNotExist(err) {
		panic(fmt.Errorf(errFileAlreadyExistTpl, outputFile))
	}
	outputFile, err = verifyFilePath(outputFile, absPath)
	if err != nil {
		panic(err)
	}

	// get fpList
	var fpList []string
	if stdin {
		// stdin
		br := bufio.NewReader(os.Stdin)
		var fp string
		for {
			fp, err = br.ReadString([]byte(delim)[0])
			if err != nil {
				if err == io.EOF {
					// normally, len(fp) equal 0
					if len(fp) > 0 {
						fp, err = verifyFilePath(fp, absPath)
						if err != nil {
							panic(err)
						}
						fpList = append(fpList, fp)
					}

					break
				}
				panic(err)
			}
			if fp == "" || fp == delim {
				continue
			}
			fp, err = verifyFilePath(fp, absPath)
			if err != nil {
				panic(err)
			}
			fpList = append(fpList, fp)
		}
	} else {
		var fp string
		fp, err = verifyFilePath(filePath, absPath)
		if err != nil {
			panic(err)
		}
		fpList = append(fpList, fp)
	}
	err = archiveTar(fpList, outputFile)
	if err != nil {
		panic(err)
	}
	os.Exit(0)
}
