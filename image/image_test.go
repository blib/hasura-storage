package image_test

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nhost/hasura-storage/image"
)

func TestManipulate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		filename string
		sum      string
		size     uint64
		options  image.Options
	}{
		{
			name:     "jpg",
			filename: "testdata/nhost.jpg",
			sum:      "12d79878831c605008cda90eb49623e969233b36f73c40646a908f6c1eb1a8e9",
			size:     33399,
			options: image.Options{
				Height:  100,
				Width:   300,
				Blur:    2,
				Quality: 50,
			},
		},
		{
			name:     "jpg",
			filename: "testdata/nhost.jpg",
			sum:      "12d79878831c605008cda90eb49623e969233b36f73c40646a908f6c1eb1a8e9",
			size:     33399,
			options:  image.Options{Width: 300, Height: 100, Blur: 2},
		},
		{
			name:     "png",
			filename: "testdata/nhost.png",
			sum:      "78cf83b463c94ecec430ca424d80650d07be5424929e67cfe978c0c753065745",
			size:     68307,
			options:  image.Options{Width: 300, Height: 100, Blur: 2},
		},
		{
			name:     "webp",
			filename: "testdata/nhost.webp",
			sum:      "01b1371cab97acc5fbadb109ddb50eeece3577ddbbeeb384ec90cd5e724c66d3",
			size:     17784,
			options:  image.Options{Width: 300, Height: 100, Blur: 2},
		},
		{
			name:     "jpg only blur",
			filename: "testdata/nhost.jpg",
			sum:      "aae488e2fdca124088dc4c3ca2daa1bb7423605d476be3c9318977044a206647",
			size:     33399,
			options:  image.Options{Blur: 2},
		},
	}

	transformer := image.NewTransformer()

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc := tc

			orig, err := os.Open(tc.filename)
			if err != nil {
				t.Fatal(err)
			}
			defer orig.Close()

			hasher := sha256.New()
			// f, _ := os.OpenFile("/tmp/nhost-test."+tc.name, os.O_WRONLY|os.O_CREATE, 0o644)
			if err := transformer.Run(orig, tc.size, hasher, tc.options); err != nil {
				t.Fatal(err)
			}

			got := hex.EncodeToString(hasher.Sum(nil))
			if !cmp.Equal(got, tc.sum) {
				t.Error(cmp.Diff(got, tc.sum))
			}
		})
	}
}

func BenchmarkManipulate(b *testing.B) {
	transformer := image.NewTransformer()
	orig, err := os.Open("testdata/nhost.jpg")
	if err != nil {
		b.Fatal(err)
	}
	defer orig.Close()
	for i := 0; i < 100; i++ {
		_, _ = orig.Seek(0, 0)

		if err := transformer.Run(
			orig,
			33399,
			io.Discard,
			image.Options{Width: 300, Height: 100, Blur: 1.5},
		); err != nil {
			b.Fatal(err)
		}
	}
}
