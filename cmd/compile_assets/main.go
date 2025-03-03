package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

func main() {
	log.SetOutput(os.Stdout)

	webPath := flag.String("web", "./web", "The absolute path for the web assets folder")
	migrationsPath := flag.String("migrations", "./migrations", "The absolute path for the migrations folder")
	outputFile := flag.String("output", "./assets/assets.bin.go", "The output Go file to dump the files into")
	flag.Parse()

	log.Println("Reading assets into memory...")

	fileMap := make(map[string][]byte)
	appendFn := func(m map[string][]byte) {
		for k, v := range m {
			fileMap[k] = v
		}
	}

	appendFn(readDir(*webPath, "web"))
	appendFn(readDir(*migrationsPath, "migrations"))

	log.Println("Writing assets go file")
	str := "package " + path.Base(path.Dir(*outputFile)) + "\n\n"
	str += "// ============================================================================\n"
	str += "// !! THIS FILE IS AUTOMATICALLY GENERATED DURING THE RELEASE/BUILD PROCESS. !!\n"
	str += "// !! You can try to overwrite it, but your changes are likely to be lost.   !!\n"
	str += "// ============================================================================\n"
	str += "\n"
	str += "// Format version: 1 (hex-encoded gzip)\n"
	str += "// Format version: 2 (base64-encoded gzip)\n"
	str += "// This file: 2\n\n"
	str += "var f2CompressedFiles = map[string]string{\n"
	for f, b := range fileMap {
		b64 := base64.StdEncoding.EncodeToString(b)
		str += fmt.Sprintf("\t\"%s\": \"%s\",\n", f, b64)
	}
	str += "}\n"
	err := os.WriteFile(*outputFile, []byte(str), 0644)
	if err != nil {
		panic(err)
	}

	log.Println("Done")
}

func readDir(dir string, pathName string) map[string][]byte {
	fileMap := make(map[string][]byte)
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fname := path.Join(dir, f.Name())
		b, err := os.ReadFile(fname)
		if err != nil {
			panic(err)
		}

		// Compress the file
		log.Println("Compressing ", fname)
		out := &bytes.Buffer{}
		gw, err := gzip.NewWriterLevel(out, gzip.BestCompression)
		if err != nil {
			panic(err)
		}
		_, _ = gw.Write(b)
		_ = gw.Close()

		fileMap[path.Join(pathName, f.Name())] = out.Bytes()
	}
	return fileMap
}
