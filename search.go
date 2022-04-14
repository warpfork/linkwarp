package linkwarp

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

/*
	Search is split up into two phases:
	first, "application" directory detection;
	second, "executable" detection, within each application directory.
	These are separated because it lets us then discuss each separately in configuration;
	and because it lets us talk about things like "the same executable name is found provided by two different applications; are you concerned?".

	Stuff like IsAppDirPredicate ended up implemented as a function, but it probably shouldn't be.
	We need rules to recognize the app dir back out of a symlink string too, for conflict resolution to be able to be stateless.
	(We could have that thing re-run pattern matching up the whole path, then reinspect, but that seems costly and not valuable.)
*/

type BinSearchCfg struct {
	SearchRoot        string                                          // Required.
	MaxDepthForAppDir int                                             // Optional.  If zero or unset, defaults to 2.
	IsAppDirPredicate func(fsys fs.FS, th string, d fs.DirEntry) bool // Optional.  A default that looks for child dirs called "bin" is used if absent.
}

func (searchCfg *BinSearchCfg) StartSearch(ctx context.Context, visitFn func(Application) error) error {
	if searchCfg.MaxDepthForAppDir == 0 {
		searchCfg.MaxDepthForAppDir = 2
	}
	if searchCfg.IsAppDirPredicate == nil {
		searchCfg.IsAppDirPredicate = func(fsys fs.FS, pth string, d fs.DirEntry) bool {
			if !d.IsDir() {
				return false
			}
			maybeBinDir, err := fsys.Open(filepath.Join(pth, "bin"))
			if err != nil {
				return false
			}
			defer maybeBinDir.Close()
			fi, err := maybeBinDir.Stat()
			if err != nil {
				return false
			}
			return fi.IsDir()
		}
	}

	searchFs := os.DirFS(searchCfg.SearchRoot)
	return fs.WalkDir(searchFs, ".", func(pth string, d fs.DirEntry, err error) error {
		if pth == "." {
			return nil
		}
		if err != nil {
			// fmt.Printf("  ?! %s\n", err) // Probably a permissions error and probably you don't care.
			return nil
		}
		isAppDir := searchCfg.IsAppDirPredicate(searchFs, pth, d)
		if isAppDir {
			fmt.Printf("  application %s -- %v\n", pth, isAppDir)
			appInfo := Application{
				Name:        filepath.Base(pth),
				Path:        pth,
				Executables: make(map[string]string),
			}
			// This inner walk is to find executables within the application dir.
			if err := fs.WalkDir(searchFs, filepath.Join(pth, "bin"), func(pth string, d fs.DirEntry, err error) error {
				// FIXME should also have a depth limit, of almost none in fact.
				fi, err := d.Info()
				if err != nil {
					return err
				}
				if strings.HasPrefix(filepath.Base(pth), ".") {
					return nil // Hidden files probably don't belong in this set.
				}
				if fi.IsDir() {
					return nil
				}
				if fi.Mode()&0111 != 0 { // "is it executable", effectively
					fmt.Printf("    bin %s\n", pth)
					appInfo.Executables[filepath.Base(pth)] = pth
				}
				return nil
			}); err != nil {
				return err
			}
			// Yield the discovery!
			if err := visitFn(appInfo); err != nil {
				return fmt.Errorf("error during update: %w", err)
			}
			return fs.SkipDir // Disallow appdirs to be seen within other appdirs.  Probably should be configurable.
		}
		// Don't search deeper if this point was already at the depth limit.
		if strings.Count(pth, string(filepath.Separator))+1 >= searchCfg.MaxDepthForAppDir {
			// fmt.Printf("  nocontinue at %s\n", pth)
			// But only return SkipDir if this is actually a dir; confusingly, the walk system skips the whole *parent* if pth isn't a dir.
			if d.IsDir() {
				return fs.SkipDir
			}
		}
		return nil
	})
}
