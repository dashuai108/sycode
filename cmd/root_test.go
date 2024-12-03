package cmd

import (
	"github.com/spf13/cobra/doc"
	"log"
	"testing"
)

func TestSyncCode(t *testing.T) {

}

func TestGenMdTree(t *testing.T) {

	err := doc.GenMarkdownTree(rootCmd, "/tmp/sycodeout")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("done!")
}
