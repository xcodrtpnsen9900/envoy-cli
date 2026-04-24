package cmd

import "sort"

// annotationCount returns the number of annotations stored for a profile.
func annotationCount(dir, profile string) (int, error) {
	annotations, err := loadAnnotations(dir, profile)
	if err != nil {
		return 0, err
	}
	return len(annotations), nil
}

// removeAnnotation deletes a single annotation key from a profile.
func removeAnnotation(dir, profile, key string) error {
	annotations, err := loadAnnotations(dir, profile)
	if err != nil {
		return err
	}
	delete(annotations, key)
	return saveAnnotations(dir, profile, annotations)
}

// sortedAnnotationKeys returns annotation keys in alphabetical order.
func sortedAnnotationKeys(annotations map[string]string) []string {
	keys := make([]string, 0, len(annotations))
	for k := range annotations {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// annotationExists reports whether a key has an annotation in the given profile.
func annotationExists(dir, profile, key string) (bool, error) {
	annotations, err := loadAnnotations(dir, profile)
	if err != nil {
		return false, err
	}
	_, ok := annotations[key]
	return ok, nil
}
