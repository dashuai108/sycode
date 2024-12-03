package sync

import (
	"archive/tar"
	"os"
	"testing"
)

func TestCopyDirToPod(t *testing.T) {
	kubeConfigPath := ""
	podName := "k8s-proxy-tenant-engine-proxy-779bf4c768-wlxvq"
	namespace := "vcluster-bc7597b9-0ad1-4e90-8f7c-788236c4e2cc"
	containerName := ""
	srcPath := "/tmp/sycodeout"
	destPath := "/tmp/test5"
	err := CopyPathToPod(kubeConfigPath, podName, namespace, containerName, srcPath, destPath)
	if err != nil {
		t.Error(err)
	}
}

func TestCopyFileToPod(t *testing.T) {
	kubeConfigPath := ""
	podName := "k8s-proxy-tenant-engine-proxy-779bf4c768-wlxvq"
	namespace := "vcluster-bc7597b9-0ad1-4e90-8f7c-788236c4e2cc"
	containerName := ""
	srcPath := "/tmp/sycodeout"
	destPath := "/tmp/test5"
	err := CopyFileToPod(kubeConfigPath, podName, namespace, containerName, srcPath, destPath)
	if err != nil {
		t.Error(err)
	}
}

func TestCompress(t *testing.T) {

	path := "/tmp/sycodeout"
	prefix := ""
	outFile, err := os.Create("example.tar")
	if err != nil {
		t.Logf("os creating tar file err: %s", err)
		t.Error(err)
	}
	defer outFile.Close()

	// 创建tar writer
	tw := tar.NewWriter(outFile)
	defer tw.Close()
	first := true

	err = compress(path, prefix, tw, first)

	if err != nil {
		t.Error(err)
	}
}
