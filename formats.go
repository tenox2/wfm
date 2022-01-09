package main

import (
	"archive/zip"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	ico "github.com/biessek/golang-ico"
	"github.com/dustin/go-humanize"
	"github.com/gen2brain/go-unarr"
	"github.com/kdomanski/iso9660"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
	"gopkg.in/ini.v1"
	"howett.net/plist"
)

func dispIcon(w http.ResponseWriter) {
	i := image.NewNRGBA(image.Rect(0, 0, 16, 16))
	d := &font.Drawer{
		Dst:  i,
		Src:  image.NewUniform(color.RGBA{0, 64, 128, 255}),
		Face: inconsolata.Bold8x16,
		Dot:  fixed.Point26_6{fixed.Int26_6(4 * 64), fixed.Int26_6(13 * 64)},
	}
	d.DrawString("W")
	w.Header().Set("Content-Type", "image/x-icon")
	ico.Encode(w, i)
}

func gourl(w http.ResponseWriter, fp string) {
	var url string
	if strings.HasSuffix(strings.ToLower(fp), ".url") {
		i, err := ini.Load(fp)
		if err != nil {
			htErr(w, "go2url", err)
			return
		}
		url = i.Section("InternetShortcut").Key("URL").String()
	}

	if strings.HasSuffix(strings.ToLower(fp), ".desktop") {
		i, err := ini.Load(fp)
		if err != nil {
			htErr(w, "go2url", err)
			return
		}
		url = i.Section("Desktop Entry").Key("URL").String()
	}

	if strings.HasSuffix(strings.ToLower(fp), ".webloc") {
		x, err := ioutil.ReadFile(fp)
		if err != nil {
			htErr(w, "go2url", err)
			return
		}
		var p struct {
			URL string
		}
		_, err = plist.Unmarshal(x, &p)
		if err != nil {
			htErr(w, "go2url", err)
			return
		}
		url = p.URL
	}

	if url == "" {
		htErr(w, "go2url", fmt.Errorf("url not found in link file"))
		return
	}
	log.Print("Redirecting to: ", url)
	redirect(w, url)
}

func listZip(w http.ResponseWriter, fp string) {
	z, err := zip.OpenReader(fp)
	if err != nil {
		htErr(w, "unzip", err)
		return
	}
	defer z.Close()
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Cache-Control", *cctl)
	for _, f := range z.File {
		fmt.Fprintf(w, "%v  %v\n", f.Name, humanize.Bytes(f.UncompressedSize64))
	}
}

func listIso(w http.ResponseWriter, fp string) {
	// TODO: recursive file list
	f, err := os.Open(fp)
	if err != nil {
		htErr(w, "isoread", err)
		return
	}
	defer f.Close()
	i, err := iso9660.OpenImage(f)
	if err != nil {
		htErr(w, "isoread", err)
		return
	}
	r, err := i.RootDir()
	if err != nil {
		htErr(w, "isoread", err)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Cache-Control", *cctl)
	if r.IsDir() {
		cld, err := r.GetChildren()
		if err != nil {
			htErr(w, "isoread", err)
			return
		}

		for _, c := range cld {
			if c.IsDir() {
				fmt.Fprintf(w, "%v  [dir]\n", c.Name())
				continue
			}
			fmt.Fprintf(w, "%v  %v\n", c.Name(), humanize.Bytes(uint64(c.Size())))
		}
	} else {
		fmt.Fprintf(w, "%v  %v\n", r.Name(), r.Size())
	}
}

func listUnarr(w http.ResponseWriter, fp string) {
	// TODO: display file sizes
	ar, err := unarr.NewArchive(fp)
	if err != nil {
		htErr(w, "unarr open", err)
		return
	}
	defer ar.Close()
	l, err := ar.List()
	if err != nil {
		htErr(w, "unarr list", err)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Cache-Control", *cctl)
	for _, e := range l {
		fmt.Fprintf(w, "%v\n", e)
	}
}
