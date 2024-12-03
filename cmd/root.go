/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/i18n"
	"os"
	"os/exec"
	"sycode/pkg/sync"
	"sycode/pkg/version"
)

type globalOptions struct {
	kubeConfigPath string // kubeconfig file path
	hostPath       string // host file path
	containerPath  string // container file path
	namespace      string // container file path
	podName        string // container file path
	containerName  string // container name
	version        string // container name
}

var (
	opts = &globalOptions{}
)

var rootCmd = &cobra.Command{
	Use:     "sycode",
	Short:   "Synchronize files to the interior of the container",
	Example: "sycode --kubeconfig=/tmp/kubeconfig --namespace=aps-os --pod=mrserver --container=test1 --host-path=/tmp/workspace/test1 --container-path=/opt/workdir/code",
	Run:     SyncCode,
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   i18n.T("Print the sycode version information"),
	Example: "sycode version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sycode version : ", version.Version)
	},
}

func SyncCode(cmd *cobra.Command, args []string) {

	// kubeConfigPath, podName, namespace, containerName, srcPath, destPath string
	logrus.Infof("kubeconfig: %s,pod-pame: %s, namespace: %s, container-name: %s,  host-path: %s, container-path: %s", opts.kubeConfigPath, opts.podName, opts.namespace, opts.containerName, opts.hostPath, opts.containerPath)
	err := sync.CopyPathToPod(opts.kubeConfigPath, opts.podName, opts.namespace, opts.containerName, opts.hostPath, opts.containerPath)
	if err != nil {
		exit(err)
	}
	logrus.Info("success!")

}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	//(default is $HOME/.cli.yaml)
	rootCmd.PersistentFlags().StringVarP(&opts.kubeConfigPath, "kubeconfig", "", "", "kubeConfig file path, The default is $HOME/.kube/config")
	rootCmd.PersistentFlags().StringVarP(&opts.namespace, "namespace", "", "", "namespace of pod (required),  Expected format is '--namespace=xxx'.")
	rootCmd.PersistentFlags().StringVarP(&opts.podName, "pod", "", "", "name of pod  (required), Expected format is '--pod=xxx'.")
	rootCmd.PersistentFlags().StringVarP(&opts.containerName, "container", "", "", "name of container.")
	rootCmd.PersistentFlags().StringVarP(&opts.hostPath, "host-path", "", "", "The host path (required), Expected format is '--host-path=/tmp/codedir'.")
	rootCmd.PersistentFlags().StringVarP(&opts.containerPath, "container-path", "", "", "The container path (required), Expected format is '--container-path=/opt/workdir/codedir'.")

	//rootCmd.MarkFlagRequired("kubeconfig")
	rootCmd.MarkPersistentFlagRequired("namespace")
	rootCmd.MarkPersistentFlagRequired("pod")
	//rootCmd.MarkFlagRequired("container")
	rootCmd.MarkPersistentFlagRequired("host-path")
	rootCmd.MarkPersistentFlagRequired("container-path")

	//rootCmd.MarkFlagsRequiredTogether("namespace", "pod", "host-path", "container-path")
}

func exit(err error) {
	var execErr *exec.ExitError
	if errors.As(err, &execErr) {
		// if there is an exit code propagate it
		exitWithCode(err, execErr.ExitCode())
	}
	// otherwise exit with catch all 1
	exitWithCode(err, 1)
}

// exits with the given error and exit code
func exitWithCode(err error, exitCode int) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(exitCode)
}
