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

	addendum "github.com/frantjc/dockerfile-addendum"
	"github.com/spf13/cobra"
)

func NewRoot() *cobra.Command {
	var (
		gz, rm, un bool
		out        string
		verbosity  int
		cmd        = &cobra.Command{
			Use:           "addendum",
			Version:       addendum.Semver(),
			SilenceErrors: true,
			SilenceUsage:  true,
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				cmd.SetContext(addendum.WithLogger(cmd.Context(), addendum.NewLogger().V(verbosity)))
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					ctx        = cmd.Context()
					logr       = addendum.LoggerFrom(ctx)
					path       = args[0]
					ext        = filepath.Ext(path)
					compressed = gz || strings.EqualFold(ext, ".tgz") || strings.EqualFold(ext, ".gz") || strings.EqualFold(ext, ".tar.gz")
					fi, err    = os.Stat(path)
				)
				switch {
				case err != nil:
					// tarball doesn't exist, so there's nothing to do
				case fi.IsDir():
					return fmt.Errorf("not a tar archive: %s", fi.Name())
				default:
					f, err := os.Open(path)
					if err != nil {
						return err
					}

					var (
						r io.Reader
					)
					if compressed {
						logr.Info("uncompressing " + f.Name())
						if r, err = gzip.NewReader(f); err != nil {
							return err
						}
					} else {
						r = f
					}

					var (
						tarball    = tar.NewReader(r)
						incomplete = true
					)
					for incomplete {
						header, err := tarball.Next()
						switch {
						case err == io.EOF:
							if rm {
								logr.Info("removing " + path)
								if err = os.Remove(path); err != nil {
									return err
								}
							}

							incomplete = false
						case err != nil:
							return err
						default:
							fullpath, err := filepath.Abs(filepath.Join(out, header.Name)) //nolint:gosec
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
					}
				}

				if un {
					if exe, err := os.Executable(); err == nil {
						if exe, err = filepath.EvalSymlinks(exe); err == nil {
							logr.Info("uninstalling " + exe)
							return os.Remove(exe)
						}

						return err
					}

					return err
				}

				return nil
			},
			Args: cobra.ExactArgs(1),
		}
	)

	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")
	cmd.Flags().BoolVarP(&rm, "rm", "r", false, "remove the tarball after extracting its contents")
	cmd.Flags().BoolVarP(&gz, "gz", "g", false, "force assuming the tarball is gzipped")
	cmd.Flags().BoolVarP(&un, "un", "u", false, "uninstall addendum on completion")
	cmd.Flags().StringVarP(&out, "out", "o", ".", "where to extract the tarball's contents to")
	cmd.Flags().CountVarP(&verbosity, "verbose", "v", "verbosity")

	return cmd
}
