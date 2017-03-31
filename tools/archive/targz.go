package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func UnTarGz(srcFilePath string, destDirPath string) error {
	fmt.Println("解压缩 " + srcFilePath + "...")
	// Create destination directory
	os.Mkdir(destDirPath, os.ModePerm)

	fr, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer fr.Close()

	// Gzip reader
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return err
	}
	defer gr.Close()

	// Tar reader
	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// End of tar archive
			break
		}
		//handleError(err)
		fmt.Printf("解压文件 %s...", hdr.Name)
		// Check if it is diretory or file
		if hdr.Typeflag != tar.TypeDir {
			// Get files from archive
			// Create diretory before create file
			err = os.MkdirAll(destDirPath, os.ModePerm)
			println(destDirPath)

			if err != nil {
				return err
			}
			// Write data to file
			fw, _ := os.Create(destDirPath + "/" + hdr.Name)
			if err != nil {
				return err
			}
			_, err = io.Copy(fw, tr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Create(root, out string) error {
	o, err := os.Create(out)
	if err != nil {
		return err
	}
	defer o.Close()
	gw := gzip.NewWriter(o)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	err = filepath.Walk(root, walkTarFunc(tw))
	if err != nil {
		return err
	}
	return nil
}

func walkTarFunc(tw *tar.Writer) filepath.WalkFunc {
	fn := func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		hdr, err := tar.FileInfoHeader(info, "")
		hdr.Name = path
		if err != nil {
			return err
		}
		err = tw.WriteHeader(hdr)
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, f)
		if err != nil {
			return err
		}
		return nil
	}
	return fn
}
