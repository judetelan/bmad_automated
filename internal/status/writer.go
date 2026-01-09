package status

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Writer writes sprint status to YAML files.
type Writer struct {
	basePath string
}

// NewWriter creates a new Writer with the specified base path.
func NewWriter(basePath string) *Writer {
	return &Writer{
		basePath: basePath,
	}
}

// UpdateStatus updates the status for a specific story key in sprint-status.yaml.
// It uses yaml.Node to preserve comments, ordering, and formatting.
func (w *Writer) UpdateStatus(storyKey string, newStatus Status) error {
	// Validate the new status
	if !newStatus.IsValid() {
		return fmt.Errorf("invalid status: %s", newStatus)
	}

	fullPath := filepath.Join(w.basePath, DefaultStatusPath)

	// Read existing file
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read sprint status: %w", err)
	}

	// Parse YAML into a Node tree to preserve formatting
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("failed to parse sprint status: %w", err)
	}

	// Find and update the story status in the node tree
	if err := updateStoryStatusInNode(&doc, storyKey, newStatus); err != nil {
		return err
	}

	// Marshal the node tree back to YAML (preserves formatting)
	updatedData, err := yaml.Marshal(&doc)
	if err != nil {
		return fmt.Errorf("failed to marshal sprint status: %w", err)
	}

	// Write back to file atomically (write to temp, then rename)
	tmpPath := fullPath + ".tmp"
	if err := os.WriteFile(tmpPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write sprint status: %w", err)
	}

	if err := os.Rename(tmpPath, fullPath); err != nil {
		// Clean up temp file on rename failure
		os.Remove(tmpPath)
		return fmt.Errorf("failed to write sprint status: %w", err)
	}

	return nil
}

// updateStoryStatusInNode finds and updates a story's status within a yaml.Node tree.
func updateStoryStatusInNode(doc *yaml.Node, storyKey string, newStatus Status) error {
	// Document node contains the root content node
	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		return fmt.Errorf("invalid YAML document structure")
	}

	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return fmt.Errorf("expected mapping at root level")
	}

	// Find development_status key in root mapping
	var devStatusNode *yaml.Node
	for i := 0; i < len(root.Content); i += 2 {
		keyNode := root.Content[i]
		if keyNode.Value == "development_status" {
			devStatusNode = root.Content[i+1]
			break
		}
	}

	if devStatusNode == nil {
		return fmt.Errorf("development_status not found in file")
	}

	if devStatusNode.Kind != yaml.MappingNode {
		return fmt.Errorf("development_status is not a mapping")
	}

	// Find the story key within development_status
	for i := 0; i < len(devStatusNode.Content); i += 2 {
		keyNode := devStatusNode.Content[i]
		if keyNode.Value == storyKey {
			// Update the value node
			valueNode := devStatusNode.Content[i+1]
			valueNode.Value = string(newStatus)
			return nil
		}
	}

	return fmt.Errorf("story not found: %s", storyKey)
}
