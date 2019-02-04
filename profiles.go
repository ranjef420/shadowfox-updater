package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-ini/ini"
	homedir "github.com/mitchellh/go-homedir"
)

func checkErr(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
		fmt.Println("You should probably report this crash in https://github.com/SrKomodo/shadowfox-updater/issues/new")
		fmt.Println("Press enter to close the program")
		fmt.Scanln()
		panic(err)
	}
}

func getProfilePaths() ([]string, []string) {
	// iniPaths stores all profiles.ini files we have to check
	iniPaths := []string{}

	// Get the home directory
	homedir, err := homedir.Dir()
	checkErr("Couldn't find home directory", err)

	// Possible places where we should check for profiles.ini
	possible := []string{
		"./profiles.ini",
		homedir + "\\AppData\\Roaming\\Mozilla\\Firefox\\profiles.ini",
		homedir + "/Library/Application Support/Firefox/profiles.ini",
		homedir + "/.mozilla/firefox/profiles.ini",
		homedir + "/.mozilla/firefox-trunk/profiles.ini",
	}

	// Check if profiles.ini exists on each possible path and add them to the list
	for _, p := range possible {
		_, err := os.Stat(p)
		if os.IsNotExist(err) {
			continue
		}
		checkErr("Couldn't check if "+p+" exists", err)
		iniPaths = append(iniPaths, p)
	}

	// If we didnt find anything then we just give up
	if len(iniPaths) == 0 {
		return nil, nil
	}

	var paths []string
	var names []string

	// For each possible ini file
	for _, p := range iniPaths {
		file, err := ini.Load(p)
		checkErr("Could not read profiles.ini, make sure its encoded in UTF-8", err)

		// Find the Path key and add it to the list
		for _, section := range file.Sections() {
			if key, err := section.GetKey("Path"); err == nil {
				path := key.String()
				isRelative := section.Key("IsRelative").MustInt(1)

				if isRelative == 1 {
					paths = append(paths, filepath.Join(filepath.Dir(p), path))
				} else {
					paths = append(paths, path)
				}
				names = append(names, filepath.Base(path))
			}
		}
	}

	return paths, names
}
