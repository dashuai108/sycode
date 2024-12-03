package sync

import (
	"archive/tar"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

func CopyFileToPod(kubeConfigPath, podName, namespace, containerName, hostPath, containerPath string) error {

	exist, err := Exist(hostPath)
	if err != nil {
		logrus.Errorf("The hostPath  [%s] not exist , error: %s !", hostPath, err)
	}
	if !exist {
		return errors.New(fmt.Sprintf("The hostPath  [%s] not exist , error: %s !", hostPath, err))
	}

	// get kubeconfig file
	defaultKubeConfig := kubeConfigFilePath(kubeConfigPath)
	if len(defaultKubeConfig) <= 0 {
		logrus.Errorf("kubeconfig not exist!")
		return errors.New("kubeconfig not exist!")
	}

	config, err := clientcmd.BuildConfigFromFlags("", defaultKubeConfig)
	if err != nil {
		logrus.Errorf("clientcmd BuildConfigFromFlags error: %s!", err)
		return err
	}

	// create k8s client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Errorf("kubernetes NewForConfig error: %s!", err)
		return err
	}

	localFile, err := os.Open(hostPath)
	defer localFile.Close()
	if err != nil {
		logrus.Errorf("os  Open error: %s!", err)
		return err
	}

	pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})

	if err != nil {
		logrus.Errorf("client get pod error: %s!", err)
		return err
	}

	var container *corev1.Container

	if len(containerName) > 0 {
		for _, c := range pod.Spec.Containers {
			if c.Name == containerName {
				container = &c
				break
			}
		}
	} else {
		container = &pod.Spec.Containers[0]
	}

	if container == nil {
		logrus.Errorf("The specified name [%s] container does not exist!", containerName)
		return errors.New(fmt.Sprintf("The specified name [%s] container does not exist!", containerName))
	}

	req := client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("container", container.Name)

	req.VersionedParams(&corev1.PodExecOptions{
		Container: container.Name,
		//Command: []string{"sh", "-c", "cat > "+destPath},
		Command: []string{"sh", "-c", fmt.Sprintf("cat > %s", containerPath)},
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())

	if err != nil {
		return err
	}

	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdin:  localFile,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    false,
	})

	if err != nil {
		return err
	}

	return nil

}

func kubeConfigFilePath(kubeConfigPath string) string {
	// get kubeconfig file
	var path string
	if len(kubeConfigPath) > 0 {
		path = kubeConfigPath
	} else {
		// get path default kubeconfig file
		path = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	return path
}

func CopyPathToPod(kubeConfigPath, podName, namespace, containerName, hostPath, containerPath string) error {

	exist, err := Exist(hostPath)
	if err != nil {
		logrus.Errorf("The hostPath  [%s] not exist , error: %s !", hostPath, err)
	}
	if !exist {
		return errors.New(fmt.Sprintf("The hostPath  [%s] not exist , error: %s !", hostPath, err))
	}

	// get kubeconfig file
	defaultKubeConfig := kubeConfigFilePath(kubeConfigPath)

	if len(defaultKubeConfig) <= 0 {
		logrus.Errorf("kubeconfig not exist!")
		return errors.New("kubeconfig not exist!")
	}

	config, err := clientcmd.BuildConfigFromFlags("", defaultKubeConfig)
	if err != nil {
		logrus.Errorf("clientcmd BuildConfigFromFlags error: %s !", err)
		return err
	}

	// create k8s client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Errorf("kubernetes NewForConfig error: %s !", err)
		return err
	}

	pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})

	if err != nil {
		logrus.Errorf("client get pod error: %s !", err)
		return err
	}

	var container *corev1.Container

	if len(containerName) > 0 {
		for _, c := range pod.Spec.Containers {
			if c.Name == containerName {
				container = &c
				break
			}
		}
	} else {
		container = &pod.Spec.Containers[0]
	}

	if container == nil {
		logrus.Errorf("The specified name [%s] container does not exist!", containerName)
		return errors.New(fmt.Sprintf("The specified name [%s] container does not exist!", containerName))
	}

	reader, writer := io.Pipe()
	tarWriter := tar.NewWriter(writer)
	go func(src string, tw *tar.Writer) {
		defer writer.CloseWithError(nil)
		defer tarWriter.Close()
		compress(hostPath, "", tarWriter, true)
	}(hostPath, tarWriter)

	req := client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("container", container.Name)

	// tar -xf mydir.tar -C /tmp/ls
	var command = []string{"sh", "-c"}
	command = append(command, fmt.Sprintf("mkdir -p %s && tar -xmf - -C %s ", containerPath, containerPath))

	req.VersionedParams(&corev1.PodExecOptions{
		Container: container.Name,
		Command:   command,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())

	if err != nil {
		logrus.Errorf("remotecommand NewSPDYExecutor error: %s ", err)
		return err
	}

	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdin:  reader,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    false,
	})

	if err != nil {
		logrus.Errorf("commit error: %s ", err)
		return err
	}

	return nil
}

func compress(path string, prefix string, tw *tar.Writer, first bool) error {
	file, _ := os.Open(path)
	info, err := file.Stat()
	if err != nil {
		logrus.Errorf("file stat error : %s ", err)
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		if first {
			prefix = ""
		}
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			logrus.Errorf("file readdir error : %s ", err)
			return err
		}
		for _, fi := range fileInfos {
			err = compress(file.Name()+"/"+fi.Name(), prefix, tw, false)
			if err != nil {
				logrus.Errorf("recursion compress error : %s ", err)
				return err
			}
		}
	} else {
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			logrus.Errorf("tar FileInfoHeader error : %s ", err)
			return err
		}
		if len(prefix) > 0 {
			header.Name = prefix[1:] + "/" + header.Name
		}

		err = tw.WriteHeader(header)
		if err != nil {
			logrus.Errorf("tar WriteHeader error : %s ", err)
			return err
		}
		_, err = io.Copy(tw, file)
		file.Close()
		if err != nil {
			logrus.Errorf("io Copy error : %s ", err)
			return err
		}
	}
	return nil
}
