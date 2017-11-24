package trimet

//go:generate msgp

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

// CSV ...
type CSV struct {
	rc io.ReadCloser `msg:"-"`
	cr *csv.Reader   `msg:"-"`
}

// Read ...
func (c *CSV) Read() ([]string, error) {
	return c.cr.Read()
}

// Close ...
func (c *CSV) Close() error {
	if c.rc == nil {
		return nil
	}
	return c.rc.Close()
}

// ReadGTFSCSV reads a GTFS txt file and returns a CSV object.
func ReadGTFSCSV(filename string) (*CSV, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	cr := csv.NewReader(f)

	cr.ReuseRecord = true
	return &CSV{cr: cr}, nil
}

// ReadZippedGTFSCSV opens a GTFS.zip file and extracts fileName as a CSV
// object.
func ReadZippedGTFSCSV(z *zip.ReadCloser, fileName string) (*CSV, error) {
	var idx int
	for i, zf := range z.File {
		if zf.Name == fileName {
			idx = i
			break
		}
	}
	rc, err := z.File[idx].Open()
	if err != nil {
		return nil, err
	}

	cr := csv.NewReader(rc)

	cr.ReuseRecord = true
	return &CSV{rc: rc, cr: cr}, nil

}

// RequestGTFSFile makes a request to download the current GTFS data set from Trimet.
// It returns an array of stops from the file.
func RequestGTFSFile(baseURL string) (*zip.ReadCloser, error) {
	f, err := ioutil.TempFile("", "tmp")
	if err != nil {
		return nil, errors.Wrap(err, "error creating tmp file")
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	resp, err := http.Get(baseURL + GTFS)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return nil, err
	}
	f.Close()

	z, err := zip.OpenReader(f.Name())
	if err != nil {
		return nil, err
	}
	return z, nil
}

func getLines(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
