package linkwarp

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"facette.io/natsort"
)

type Application struct {
	Name        string            // Name
	Path        string            // Path is the whole path, relative to the searchroot.
	Executables map[string]string // Executables is a map from short name to full path (again, relative to searchroot) of executables known in this application.
}

type BinSynthesisCfg struct {
	AppRoot string
	BinRoot string
}

func (synthCfg *BinSynthesisCfg) UpdateLinks(appInfo Application) error {
	// This is not handling race conditions particularly calmly, but neither is it doing anything dangerously wrong in those cases.
	// Rerunning the system should cause it to converge.
	for name, pth := range appInfo.Executables {
		if err := func() error {

			lnkPath := filepath.Join(synthCfg.BinRoot, name)
			targetPath := filepath.Join(synthCfg.AppRoot, pth)

			fi, err := os.Lstat(lnkPath)
			if err != nil {
				if os.IsNotExist(err) {
					// Great.  We get to make it.
					if err := os.Symlink(targetPath, lnkPath); err != nil {
						return fmt.Errorf("creating link failed: %w", err)
					}
					return nil
				} else {
					return err
				}
			}
			// If something exists and it's a symlink, read it, and see if it "wins"; leave it be if so; otherwise replace it.
			if fi.Mode()&os.ModeSymlink != 0 {
				existing, err := os.Readlink(lnkPath)
				if err != nil {
					return fmt.Errorf("can't read existing link, so not sure how to procede: %w", err)
				}
				if existing == targetPath {
					return nil
				}
				// Compare.
				// Re-infer the application name from the symlink's target string, then size it up with the new one.
				existingChunks := strings.Split(existing, string(filepath.Separator))
				existingAppName := ""
				for i := len(existingChunks) - 1; i > 0; i-- {
					if existingChunks[i] == "bin" { // note that this is the bit that makes a fool out of IsAppDirPredicate.  hard to be that flexible here.
						existingAppName = existingChunks[i-1]
						break
					}
				}
				// If the existing one looks bigger, skip out with no changes.
				if !natsort.Compare(existingAppName, appInfo.Name) {
					fmt.Printf("         leaving link for %q: %q looked older than %q\n", name, appInfo.Name, existingAppName)
					return nil
				}
				// Overwrite.
				//  Creating a temp file and renaming it into place.  This is borderline overkill, but it if you want a zero-downtime in-place-updates behavior, this is what you want.
				fmt.Printf("         updating link for %q: %q looked newer than %q\n", name, appInfo.Name, existingAppName)
				for {
					tmpName := strconv.Itoa(int(time.Now().Unix()))
					tmpName = filepath.Join(synthCfg.BinRoot, ".tmp."+tmpName)
					if err := os.Symlink(targetPath, tmpName); err != nil {
						if os.IsExist(err) {
							continue
						}
						return fmt.Errorf("updating link failed: %w", err)
					}
					return os.Rename(tmpName, lnkPath)
				}
				return nil
			}
			// If something existed and it wasn't a symlink, we're... not gonna touch that.
			return fmt.Errorf("something existed at %q and wasn't a symlink; exiting in fear", lnkPath)
		}(); err != nil {
			return err
		}
	}
	return nil
}
