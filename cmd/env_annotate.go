package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func annotationFilePath(dir, profile string) string {
	return filepath.Join(dir, ".envoy", "annotations", profile+".annotations")
}

func init() {
	annotateCmd := &cobra.Command{
		Use:   "annotate",
		Short: "Manage inline annotations (comments) for profile keys",
	}

	setAnnotation := &cobra.Command{
		Use:   "set <profile> <key> <note>",
		Short: "Set an annotation for a key in a profile",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			dir := projectDir()
			if err := setAnnotation(dir, args[0], args[1], args[2]); err != nil {
				fatalf("annotate set: %v", err)
			}
			fmt.Printf("Annotation set for key %q in profile %q\n", args[1], args[0])
		},
	}

	getAnnotation := &cobra.Command{
		Use:   "get <profile> <key>",
		Short: "Get the annotation for a key in a profile",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			dir := projectDir()
			note, err := getAnnotationValue(dir, args[0], args[1])
			if err != nil {
				fatalf("annotate get: %v", err)
			}
			if note == "" {
				fmt.Printf("No annotation for key %q\n", args[1])
			} else {
				fmt.Printf("%s = %s\n", args[1], note)
			}
		},
	}

	listAnnotations := &cobra.Command{
		Use:   "list <profile>",
		Short: "List all annotations for a profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dir := projectDir()
			annotations, err := loadAnnotations(dir, args[0])
			if err != nil {
				fatalf("annotate list: %v", err)
			}
			if len(annotations) == 0 {
				fmt.Println("No annotations found.")
				return
			}
			for k, v := range annotations {
				fmt.Printf("  %s: %s\n", k, v)
			}
		},
	}

	annotateCmd.AddCommand(setAnnotation, getAnnotation, listAnnotations)
	rootCmd.AddCommand(annotateCmd)
}

func setAnnotation(dir, profile, key, note string) error {
	annotations, _ := loadAnnotations(dir, profile)
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[key] = note
	return saveAnnotations(dir, profile, annotations)
}

func getAnnotationValue(dir, profile, key string) (string, error) {
	annotations, err := loadAnnotations(dir, profile)
	if err != nil {
		return "", err
	}
	return annotations[key], nil
}

func loadAnnotations(dir, profile string) (map[string]string, error) {
	path := annotationFilePath(dir, profile)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, err
	}
	annotations := map[string]string{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			annotations[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return annotations, nil
}

func saveAnnotations(dir, profile string, annotations map[string]string) error {
	path := annotationFilePath(dir, profile)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	var sb strings.Builder
	for k, v := range annotations {
		sb.WriteString(k + "=" + v + "\n")
	}
	return os.WriteFile(path, []byte(sb.String()), 0644)
}
