package storage

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

type StorageService interface {
	Upload(file *multipart.FileHeader, filename string) (string, error)
	Delete(fileUrl string) error
}

type supabaseStorage struct {
	url        string
	serviceKey string
	bucket     string
}

func NewSupabaseStorage(url, serviceKey, bucket string) StorageService {
	// Ensure URL doesn't end with slash
	url = strings.TrimRight(url, "/")
	return &supabaseStorage{
		url:        url,
		serviceKey: serviceKey,
		bucket:     bucket,
	}
}

func (s *supabaseStorage) Upload(file *multipart.FileHeader, filename string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Read file content
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return "", err
	}

	// Detect content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		// Fallback detection if header is missing
		contentType = http.DetectContentType(fileBytes)
	}

	// Prepare request
	// URL format: https://[project_ref].supabase.co/storage/v1/object/[bucket]/[path]
	apiPath := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.url, s.bucket, filename)

	req, err := http.NewRequest("POST", apiPath, bytes.NewBuffer(fileBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+s.serviceKey)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to upload to supabase: status %d, body: %s", resp.StatusCode, string(body))
	}

	// Construct public URL
	// Public URL format: https://[project_ref].supabase.co/storage/v1/object/public/[bucket]/[path]
	publicUrl := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.url, s.bucket, filename)
	return publicUrl, nil
}

func (s *supabaseStorage) Delete(fileUrl string) error {
	// We need to extract the path relative to the bucket.
	// Expected fileUrl: https://[project_ref].supabase.co/storage/v1/object/public/[bucket]/[file_path]
	// API endpoint for delete: DELETE https://[project_ref].supabase.co/storage/v1/object/[bucket]/[file_path]

	// Simply replace "/public/" with "/" to get the private object path for deletion?
	// The delete API is: DELETE /storage/v1/object/{bucket}/{wildcard}

	// Let's parse the URL to extract the key.
	// A robust way: find "/storage/v1/object/public/[bucket]/" and take everything after.

	pattern := fmt.Sprintf("/storage/v1/object/public/%s/", s.bucket)
	idx := strings.Index(fileUrl, pattern)
	if idx == -1 {
		// If it doesn't match our pattern, maybe it's not hosted here or format changed.
		// Fail gracefully or log warning?
		return fmt.Errorf("invalid file url for deletion: %s", fileUrl)
	}

	// Extract 'filename.ext'
	filePath := fileUrl[idx+len(pattern):]

	// API URL
	apiPath := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.url, s.bucket, filePath)

	req, err := http.NewRequest("DELETE", apiPath, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.serviceKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete from supabase: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Helper to check if file is an image - keeping logic from service if needed,
// but service already validates.
func IsImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}
