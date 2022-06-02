package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/frantjc/dockerfile-addendum"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "addendum",
		Version: addendum.Semver(),
		RunE:    run,
		Args:    cobra.ExactArgs(1),
	}
	rm  bool
	out string
	gz  bool
)

func init() {
	rootCmd.SetVersionTemplate(
		fmt.Sprintf("{{ .Name }}{{ .Version }} %s\n", runtime.Version()),
	)
	rootCmd.Flags().BoolVar(&rm, "rm", false, "Remove the tarball after extracting its contents")
	rootCmd.Flags().StringVarP(&out, "out", "o", ".", "Where to extract the tarball's contents to")
	rootCmd.Flags().BoolVar(&gz, "gz", false, "Force assuming the tarball is gzipped")
}

func run(cmd *cobra.Command, args []string) error {
	var (
		path       = args[0]
		ext        = filepath.Ext(path)
		compressed = gz || strings.EqualFold(ext, ".tgz") || strings.EqualFold(ext, ".gz") || strings.EqualFold(ext, ".tar.gz")
		r          io.Reader
		fi, err    = os.Stat(path)
	)
	switch {
	case err != nil:
		// tarball doesn't exist, so there's nothing to do
		return nil
	case fi.IsDir():
		return fmt.Errorf("directory %s is not a tar archive", fi.Name())
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	if compressed {
		r, err = gzip.NewReader(f)
	} else {
		r = f
	}

	tarball := tar.NewReader(r)

	for {
		header, err := tarball.Next()
		switch {
		case err == io.EOF:
			if rm {
				return os.Remove(path)
			}

			return nil
		case err != nil:
			return err
		}

		fullpath, err := filepath.Abs(
			filepath.Join(out, header.Name),
		)
		if err != nil {
			return fmt.Errorf("unable to determine path for tar header %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			di, err := os.Stat(fullpath)
			switch {
			case err == nil && !di.IsDir():
				return fmt.Errorf("not a directory %s", fullpath)
			case err == nil && di.IsDir():
				// nothing to do
			default:
				if err := os.Mkdir(fullpath, header.FileInfo().Mode().Perm()); err != nil {
					return fmt.Errorf("unable to create directory %s", fullpath)
				}
			}
		case tar.TypeReg:
			o, err := os.Create(fullpath)
			if err != nil {
				return fmt.Errorf("unable to create file %s", fullpath)
			}
			defer o.Close()

			if _, err := io.CopyN(o, tarball, header.Size); err != nil {
				return fmt.Errorf("unable to write to file %s", fullpath)
			}
		default:
			return fmt.Errorf("unable to handle tar header type %b", header.Typeflag)
		}
	}
}

func main() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
