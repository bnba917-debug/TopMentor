package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrInvalidKind    = fmt.Errorf("invalid upload kind")
	ErrFileTooLarge   = fmt.Errorf("file too large")
	ErrUnsupportedExt = fmt.Errorf("unsupported file type")
)

const (
	KindAvatar       = "avatar"
	KindIntroVideo   = "intro_video"
	KindIDCard       = "id_card"
	KindStudentCard  = "student_card"
	KindEnglishProof = "english_proof"
)

type Store struct {
	rootDir string
}

func NewStore(rootDir string) (*Store, error) {
	if rootDir == "" {
		rootDir = "uploads"
	}
	if err := os.MkdirAll(rootDir, 0o755); err != nil {
		return nil, err
	}
	return &Store{rootDir: rootDir}, nil
}

func (s *Store) SaveMentorFile(mentorID int64, kind string, fh *multipart.FileHeader) (string, error) {
	maxBytes, allowedExts := limitsForKind(kind)
	if maxBytes == 0 {
		return "", ErrInvalidKind
	}
	if fh.Size > maxBytes {
		return "", ErrFileTooLarge
	}

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !allowedExts[ext] {
		return "", ErrUnsupportedExt
	}

	dir := filepath.Join(s.rootDir, "mentors", fmt.Sprintf("%d", mentorID))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	name := fmt.Sprintf("%s_%d%s", kind, time.Now().UnixNano(), ext)
	dest := filepath.Join(dir, name)

	src, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, io.LimitReader(src, maxBytes+1)); err != nil {
		return "", err
	}

	return "/uploads/mentors/" + fmt.Sprintf("%d", mentorID) + "/" + name, nil
}

func (s *Store) SaveApplicantFile(userID int64, kind string, fh *multipart.FileHeader) (string, error) {
	maxBytes, allowedExts := limitsForKind(kind)
	if maxBytes == 0 {
		return "", ErrInvalidKind
	}
	if fh.Size > maxBytes {
		return "", ErrFileTooLarge
	}

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !allowedExts[ext] {
		return "", ErrUnsupportedExt
	}

	dir := filepath.Join(s.rootDir, "applicants", fmt.Sprintf("%d", userID))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	name := fmt.Sprintf("%s_%d%s", kind, time.Now().UnixNano(), ext)
	dest := filepath.Join(dir, name)

	src, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, io.LimitReader(src, maxBytes+1)); err != nil {
		return "", err
	}

	return "/uploads/applicants/" + fmt.Sprintf("%d", userID) + "/" + name, nil
}

func limitsForKind(kind string) (int64, map[string]bool) {
	docExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".pdf": true}
	switch kind {
	case KindAvatar:
		return 2 << 20, map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	case KindIntroVideo:
		return 30 << 20, map[string]bool{".mp4": true, ".webm": true, ".mov": true}
	case KindIDCard, KindStudentCard, KindEnglishProof:
		return 5 << 20, docExts
	default:
		return 0, nil
	}
}
