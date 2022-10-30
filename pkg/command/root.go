package command

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/frantjc/dockerfile-addendum"
	"github.com/spf13/cobra"
)

func NewRoot() *cobra.Command {
	var (
		gz, rm, un bool
		out        string
		cmd        = &cobra.Command{
			Use:     "addendum",
			Version: addendum.Semver(),
			RunE: func(cmd *cobra.Command, args []string) error {
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
					return fmt.Errorf("not a tar archive: %s", fi.Name())
				}

				f, err := os.Open(path)
				if err != nil {
					return err
				}

				if compressed {
					if r, err = gzip.NewReader(f); err != nil {
						return err
					}
				} else {
					r = f
				}

				tarball := tar.NewReader(r)

				for {
					header, err := tarball.Next()
					switch {
					case err == io.EOF:
						if rm {
							if err = os.Remove(path); err != nil {
								return err
							}
						}

						if un {
							if exe, err := os.Executable(); err == nil {
								return os.Remove(exe)
							} else {
								return err
							}
						}

						return nil
					case err != nil:
						return err
					}

					fullpath, err := filepath.Abs(
						filepath.Join(out, header.Name),
					)
					if err != nil {
						return fmt.Errorf("determine path for tar header: %s", header.Name)
					}

					switch header.Typeflag {
					case tar.TypeDir:
						di, err := os.Stat(fullpath)
						switch {
						case err == nil && !di.IsDir():
							return fmt.Errorf("not a directory: %s", fullpath)
						case err == nil && di.IsDir():
							// nothing to do
						default:
							if err := os.Mkdir(fullpath, header.FileInfo().Mode().Perm()); err != nil {
								return fmt.Errorf("create directory: %s", fullpath)
							}
						}
					case tar.TypeReg:
						o, err := os.Create(fullpath)
						if err != nil {
							return fmt.Errorf("create file: %s", fullpath)
						}
						defer o.Close()

						if _, err := io.CopyN(o, tarball, header.Size); err != nil {
							return fmt.Errorf("write to file: %s", fullpath)
						}
					default:
						return fmt.Errorf("handle tar header type: %b", header.Typeflag)
					}
				}
			},
			Args: cobra.ExactArgs(1),
		}
	)

	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")
	cmd.Flags().BoolVarP(&rm, "rm", "r", false, "Remove the tarball after extracting its contents")
	cmd.Flags().BoolVarP(&gz, "gz", "g", false, "Force assuming the tarball is gzipped")
	cmd.Flags().BoolVarP(&un, "un", "u", false, "Uninstall addendum on completion")
	cmd.Flags().StringVarP(&out, "out", "o", ".", "Where to extract the tarball's contents to")

	return cmd
}
