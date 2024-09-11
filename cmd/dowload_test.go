package cmd

import (
	"fmt"
	"testing"

	"github.com/cisco-open/grabit/test"
	"github.com/stretchr/testify/assert"
)

func TestRunDownload(t *testing.T) {
	content := `abcdef`
	port := test.TestHttpHandler(content, t)
	testfilepath := test.TmpFile(t, fmt.Sprintf(`
	[[Resource]]
	Urls = ['http://localhost:%d/test.html']
	Integrity = 'sha256-vvV+x/U6bUC+tkCngKY5yDvCmsipgW8fxsXG3Nk8RyE='

	[[Resource]]
	Urls = ['http://localhost:%d/test3.html']
	Integrity = 'sha256-vvV+x/U6bUC+tkCngKY5yDvCmsipgW8fxsXG3Nk8RyE='
`, port, port))
	outputDir := test.TmpDir(t)
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"-f", testfilepath, "download", "--dir", outputDir})
	err := cmd.Execute()
	assert.Nil(t, err)
	for _, file := range []string{"test.html", "test3.html"} {
		test.AssertFileContains(t, fmt.Sprintf("%s/%s", outputDir, file), content)
	}
}

func TestRunDownloadWithTags(t *testing.T) {
	content := `abcdef`
	port := test.TestHttpHandler(content, t)
	testfilepath := test.TmpFile(t, fmt.Sprintf(`
	[[Resource]]
	Urls = ['http://localhost:%d/test.html']
	Integrity = 'sha256-vvV+x/U6bUC+tkCngKY5yDvCmsipgW8fxsXG3Nk8RyE='
	Tags = ['tag']

	[[Resource]]
	Urls = ['http://localhost:%d/test2.html']
	Integrity = 'sha256-vvV+x/U6bUC+tkCngKY5yDvCmsipgW8fxsXG3Nk8RyE='
	Tags = ['tag1', 'tag2']
`, port, port))
	outputDir := test.TmpDir(t)
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"-f", testfilepath, "download", "--tag", "tag", "--dir", outputDir})
	err := cmd.Execute()
	assert.Nil(t, err)
	for _, file := range []string{"test.html"} {
		test.AssertFileContains(t, fmt.Sprintf("%s/%s", outputDir, file), content)
	}
}

func TestRunDownloadWithoutTags(t *testing.T) {
	content := `abcdef`
	port := test.TestHttpHandler(content, t)
	testfilepath := test.TmpFile(t, fmt.Sprintf(`
	[[Resource]]
	Urls = ['http://localhost:%d/test.html']
	Integrity = 'sha256-vvV+x/U6bUC+tkCngKY5yDvCmsipgW8fxsXG3Nk8RyE='
	Tags = ['tag']

	[[Resource]]
	Urls = ['http://localhost:%d/test2.html']
	Integrity = 'sha256-vvV+x/U6bUC+tkCngKY5yDvCmsipgW8fxsXG3Nk8RyE='
	Tags = ['tag1', 'tag2']
`, port, port))
	outputDir := test.TmpDir(t)
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"-f", testfilepath, "download", "--notag", "tag", "--dir", outputDir})
	err := cmd.Execute()
	assert.Nil(t, err)
	for _, file := range []string{"test2.html"} {
		test.AssertFileContains(t, fmt.Sprintf("%s/%s", outputDir, file), content)
	}
}

func TestRunDownloadMultipleErrors(t *testing.T) {
	testfilepath := test.TmpFile(t, `
	[[Resource]]
	Urls = ['http://localhost:1234/test.html']
	Integrity = 'sha256-vvV+x/U6bUC+tkCngKY5yDvCmsipgW8fxsXG3Nk8RyE='

	[[Resource]]
	Urls = ['http://cannot-be-resolved.no:12/test.html']
	Integrity = 'sha256-vvV+x/U6bUC+tkCngKY5yDvCmsipgW8fxsXG3Nk8RyE='
`)
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"-f", testfilepath, "download"})
	err := cmd.Execute()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to download")
	assert.Contains(t, err.Error(), "connection refused")
	assert.Contains(t, err.Error(), "no such host")
}
