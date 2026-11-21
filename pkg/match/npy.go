package match

import (
	"fmt"
	"os"
)

// LoadPatchBytes reads patch bytes from a .npy file.
// Format: numpy array of shape (N, 128) with dtype uint8.
func LoadPatchBytes(path string) ([]byte, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot open %s: %w", path, err)
	}
	defer f.Close()

	// Read npy header: magic(6) + major(1) + minor(1) + header_len(2) + header_str
	header := make([]byte, 6)
	if _, err := f.Read(header); err != nil {
		return nil, 0, fmt.Errorf("read npy magic: %w", err)
	}
	if string(header[:6]) != "\x93NUMPY" {
		return nil, 0, fmt.Errorf("not a valid npy file")
	}

	// Read version
	ver := make([]byte, 2)
	if _, err := f.Read(ver); err != nil {
		return nil, 0, err
	}

	// Read header length (2 bytes for v1, 4 bytes for v2+)
	var headerLen int
	if ver[0] == 1 {
		bl := make([]byte, 2)
		f.Read(bl)
		headerLen = int(bl[0]) | int(bl[1])<<8
	} else {
		bl := make([]byte, 4)
		f.Read(bl)
		headerLen = int(bl[0]) | int(bl[1])<<8 | int(bl[2])<<16 | int(bl[3])<<24
	}

	// Read header string
	headerStr := make([]byte, headerLen)
	if _, err := f.Read(headerStr); err != nil {
		return nil, 0, err
	}

	// Parse shape from header: "('uint8', False, (N, 128))"
	hdr := string(headerStr)
	shapeStart := -1
	for i := len(hdr) - 1; i >= 0; i-- {
		if hdr[i] == '(' {
			shapeStart = i
			break
		}
	}
	if shapeStart < 0 {
		return nil, 0, fmt.Errorf("cannot find shape in npy header")
	}
	shapeStr := hdr[shapeStart:]
	n1, n2 := 0, 0
	fmt.Sscanf(shapeStr, "(%d, %d)", &n1, &n2)
	if n1 == 0 || n2 != 128 {
		return nil, 0, fmt.Errorf("unexpected shape %d,%d in npy", n1, n2)
	}

	// Read data
	data := make([]byte, n1*128)
	if _, err := f.Read(data); err != nil {
		return nil, 0, err
	}

	return data, n1, nil
}
