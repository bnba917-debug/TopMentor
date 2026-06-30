package upload

import (
	"bytes"
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_SaveMentorFile_Avatar(t *testing.T) {
	dir := t.TempDir()
	store, err := NewStore(dir)
	require.NoError(t, err)

	fh := fakeFileHeader("photo.jpg", "image/jpeg", []byte("fake-image"))
	url, err := store.SaveMentorFile(7, KindAvatar, fh)
	require.NoError(t, err)
	assert.Contains(t, url, "/uploads/mentors/7/avatar_")
	assert.Contains(t, url, ".jpg")
}

func TestStore_SaveMentorFile_RejectLarge(t *testing.T) {
	dir := t.TempDir()
	store, err := NewStore(dir)
	require.NoError(t, err)

	big := make([]byte, 3<<20)
	fh := fakeFileHeader("big.jpg", "image/jpeg", big)
	_, err = store.SaveMentorFile(1, KindAvatar, fh)
	assert.ErrorIs(t, err, ErrFileTooLarge)
}

func fakeFileHeader(name, contentType string, body []byte) *multipart.FileHeader {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, _ := w.CreatePart(textproto.MIMEHeader{
		"Content-Disposition": []string{`form-data; name="file"; filename="` + name + `"`},
		"Content-Type":        []string{contentType},
	})
	_, _ = part.Write(body)
	_ = w.Close()

	r := multipart.NewReader(&buf, w.Boundary())
	form, _ := r.ReadForm(10 << 20)
	return form.File["file"][0]
}
