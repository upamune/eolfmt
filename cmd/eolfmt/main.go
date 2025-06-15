// Package main implements the eolfmt command-line tool.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"

	"github.com/sourcegraph/conc/pool"
	"github.com/upamune/eolfmt/internal/processor"
	"github.com/upamune/eolfmt/internal/walker"
)

// Build-time variables injected via -ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type stats struct {
	processed atomic.Uint64
	modified  atomic.Uint64
	errors    atomic.Uint64
}

func main() {
	var versionFlag bool
	flag.BoolVar(&versionFlag, "version", false, "show version information")
	flag.Parse()

	if versionFlag {
		fmt.Printf("eolfmt version %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  built at: %s\n", date)
		fmt.Printf("  go version: %s\n", runtime.Version())
		os.Exit(0)
	}

	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx, paths); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, paths []string) error {
	var stats stats

	maxWorkers := runtime.GOMAXPROCS(0)
	p := pool.New().
		WithMaxGoroutines(maxWorkers).
		WithContext(ctx)

	for _, path := range paths {
		ch := walker.WalkFiles(path)
		for {
			select {
			case <-ctx.Done():
				printStats(&stats, true)
				return ctx.Err()
			case fileInfo, ok := <-ch:
				if !ok {
					goto nextPath
				}
				fi := fileInfo
				p.Go(func(_ context.Context) error {
					stats.processed.Add(1)

					modified, err := processor.CheckAndFixFile(fi.Path, fi.Info)
					if err != nil {
						if !os.IsPermission(err) {
							stats.errors.Add(1)
						}
					} else if modified {
						stats.modified.Add(1)
					}
					return nil
				})
			}
		}
	nextPath:
	}

	if err := p.Wait(); err != nil {
		return err
	}

	printStats(&stats, false)
	return nil
}

func printStats(s *stats, interrupted bool) {
	if interrupted {
		fmt.Print("eolfmt: interrupted - ")
	} else {
		fmt.Print("eolfmt: ")
	}
	fmt.Printf("%d files processed, %d modified",
		s.processed.Load(), s.modified.Load())
	if e := s.errors.Load(); e > 0 {
		fmt.Printf(", %d errors", e)
	}
	fmt.Println()
}
