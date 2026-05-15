package service

import (
	"encoding/json"
	"image"

	maa "github.com/MaaXYZ/maa-framework-go/v4"
)

func serializeTaskDetail(detail *maa.TaskDetail) map[string]interface{} {
	if detail == nil {
		return nil
	}

	nodes := make([]map[string]interface{}, 0, len(detail.Nodes))
	for _, nodeRef := range detail.Nodes {
		nodeDetail, err := nodeRef.GetDetail()
		if err != nil {
			nodes = append(nodes, map[string]interface{}{
				"id":    nodeRef.ID(),
				"error": err.Error(),
			})
			continue
		}
		nodes = append(nodes, serializeNodeDetail(nodeDetail))
	}

	return map[string]interface{}{
		"id":     detail.ID,
		"entry":  detail.Entry,
		"status": detail.Status.String(),
		"nodes":  nodes,
	}
}

func serializeNodeDetail(detail *maa.NodeDetail) map[string]interface{} {
	if detail == nil {
		return nil
	}

	result := map[string]interface{}{
		"id":            detail.ID,
		"name":          detail.Name,
		"runCompleted":  detail.RunCompleted,
		"run_completed": detail.RunCompleted,
	}
	if detail.Recognition != nil {
		result["recognition"] = serializeRecognitionDetail(detail.Recognition)
	}
	if detail.Action != nil {
		result["action"] = serializeActionDetail(detail.Action)
	}
	return result
}

func serializeRecognitionDetail(detail *maa.RecognitionDetail) map[string]interface{} {
	if detail == nil {
		return nil
	}

	result := map[string]interface{}{
		"id":        detail.ID,
		"name":      detail.Name,
		"algorithm": detail.Algorithm,
		"hit":       detail.Hit,
		"box":       rectToSlice(detail.Box),
	}
	if parsed, ok := parseDetailJSON(detail.DetailJson); ok {
		result["detail"] = parsed
	} else if detail.DetailJson != "" {
		result["detail_json"] = detail.DetailJson
	}
	if detail.Results != nil {
		result["results"] = serializeRecognitionResults(detail.Results)
	}
	if len(detail.CombinedResult) > 0 {
		combined := make([]map[string]interface{}, 0, len(detail.CombinedResult))
		for _, child := range detail.CombinedResult {
			combined = append(combined, serializeRecognitionDetail(child))
		}
		result["combined_result"] = combined
	}
	if detail.Raw != nil {
		result["raw_image"] = imageMeta(detail.Raw)
	}
	if len(detail.Draws) > 0 {
		draws := make([]map[string]interface{}, 0, len(detail.Draws))
		for _, draw := range detail.Draws {
			draws = append(draws, imageMeta(draw))
		}
		result["draws"] = draws
	}
	return result
}

func serializeActionDetail(detail *maa.ActionDetail) map[string]interface{} {
	if detail == nil {
		return nil
	}

	result := map[string]interface{}{
		"id":      detail.ID,
		"name":    detail.Name,
		"action":  detail.Action,
		"box":     rectToSlice(detail.Box),
		"success": detail.Success,
	}
	if parsed, ok := parseDetailJSON(detail.DetailJson); ok {
		result["detail"] = parsed
	} else if detail.DetailJson != "" {
		result["detail_json"] = detail.DetailJson
	}
	if detail.Result != nil {
		result["result_type"] = string(detail.Result.Type())
		result["result"] = detail.Result.Value()
	}
	return result
}

func serializeRecognitionResults(results *maa.RecognitionResults) map[string]interface{} {
	if results == nil {
		return nil
	}

	return map[string]interface{}{
		"all":      serializeRecognitionResultList(results.All),
		"best":     serializeRecognitionResult(results.Best),
		"filtered": serializeRecognitionResultList(results.Filtered),
	}
}

func serializeRecognitionResultList(items []*maa.RecognitionResult) []map[string]interface{} {
	output := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		output = append(output, serializeRecognitionResult(item))
	}
	return output
}

func serializeRecognitionResult(result *maa.RecognitionResult) map[string]interface{} {
	if result == nil {
		return nil
	}

	output := map[string]interface{}{
		"type": string(result.Type()),
	}
	if value := result.Value(); value != nil {
		output["value"] = value
	}
	return output
}

func rectToSlice(rect maa.Rect) []int {
	return []int{rect.X(), rect.Y(), rect.Width(), rect.Height()}
}

func parseDetailJSON(detail string) (interface{}, bool) {
	if detail == "" {
		return nil, false
	}
	var parsed interface{}
	if err := json.Unmarshal([]byte(detail), &parsed); err != nil {
		return nil, false
	}
	return parsed, true
}

func imageMeta(img image.Image) map[string]interface{} {
	if img == nil {
		return nil
	}
	bounds := img.Bounds()
	return map[string]interface{}{
		"width":  bounds.Dx(),
		"height": bounds.Dy(),
	}
}
