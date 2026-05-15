package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func resolvePipelineEntryFile(pipelineName string) (string, error) {
	if pipelineName == "" {
		return "", fmt.Errorf("任务名不能为空")
	}

	filePath := filepath.Join("resource", "pipeline", pipelineName+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取 Pipeline 文件失败: %w", err)
	}

	var nodes map[string]interface{}
	if err := json.Unmarshal(data, &nodes); err != nil {
		return "", fmt.Errorf("解析 Pipeline 文件失败: %w", err)
	}

	return resolvePipelineEntry(pipelineName, nodes)
}

func ResolvePipelineEntryName(pipelineName string) (string, error) {
	return resolvePipelineEntryFile(pipelineName)
}

func resolvePipelineEntry(pipelineName string, nodes map[string]interface{}) (string, error) {
	if _, ok := nodes[pipelineName]; ok {
		return pipelineName, nil
	}

	roots := findPipelineRootNodes(nodes)
	if len(roots) == 1 {
		return roots[0], nil
	}
	if len(roots) == 0 {
		return "", fmt.Errorf("Pipeline %q 没有可自动判断的入口节点，请创建同名节点或手动执行指定节点", pipelineName)
	}
	return "", fmt.Errorf("Pipeline %q 有多个入口节点 %s，请创建同名节点或在编辑器中执行指定节点", pipelineName, strings.Join(roots, ", "))
}

func findPipelineRootNodes(nodes map[string]interface{}) []string {
	children := map[string]bool{}
	for _, rawNode := range nodes {
		node, ok := rawNode.(map[string]interface{})
		if !ok {
			continue
		}
		for _, target := range relationTargetNames(node["next"]) {
			children[target] = true
		}
	}

	roots := make([]string, 0, len(nodes))
	for name := range nodes {
		if !children[name] {
			roots = append(roots, name)
		}
	}
	sort.Strings(roots)
	return roots
}

func relationTargetNames(value interface{}) []string {
	var result []string
	for _, item := range relationItems(value) {
		switch typed := item.(type) {
		case string:
			target := concreteRelationName(typed)
			if target != "" {
				result = append(result, target)
			}
		case map[string]interface{}:
			if anchor, _ := typed["anchor"].(bool); anchor {
				continue
			}
			if name, _ := typed["name"].(string); strings.TrimSpace(name) != "" {
				result = append(result, strings.TrimSpace(name))
			}
		}
	}
	return result
}

func relationItems(value interface{}) []interface{} {
	if value == nil {
		return nil
	}
	if items, ok := value.([]interface{}); ok {
		return items
	}
	return []interface{}{value}
}

func concreteRelationName(value string) string {
	text := strings.TrimSpace(value)
	if text == "" || strings.HasPrefix(text, "[Anchor]") {
		return ""
	}
	if strings.HasPrefix(text, "[JumpBack]") {
		return strings.TrimSpace(strings.TrimPrefix(text, "[JumpBack]"))
	}
	return text
}
