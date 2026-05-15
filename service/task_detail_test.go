package service

import (
	"testing"

	maa "github.com/MaaXYZ/maa-framework-go/v4"
)

func TestSerializeRecognitionDetailIncludesCoreDebugFields(t *testing.T) {
	detail := &maa.RecognitionDetail{
		ID:         7,
		Name:       "ReadTitle",
		Algorithm:  "OCR",
		Hit:        true,
		Box:        maa.Rect{12, 34, 56, 78},
		DetailJson: `{"all":[{"box":[12,34,56,78],"text":"OK","score":0.98}]}`,
	}

	got := serializeRecognitionDetail(detail)

	if got["id"] != int64(7) {
		t.Fatalf("expected id 7, got %#v", got["id"])
	}
	if got["name"] != "ReadTitle" || got["algorithm"] != "OCR" || got["hit"] != true {
		t.Fatalf("unexpected recognition summary: %#v", got)
	}
	if box, ok := got["box"].([]int); !ok || len(box) != 4 || box[2] != 56 {
		t.Fatalf("expected serialized box, got %#v", got["box"])
	}
	if detailJSON, ok := got["detail"].(map[string]interface{}); !ok || detailJSON["all"] == nil {
		t.Fatalf("expected parsed detail JSON, got %#v", got["detail"])
	}
}

func TestSerializeActionDetailIncludesParsedDetail(t *testing.T) {
	detail := &maa.ActionDetail{
		ID:         9,
		Name:       "TapConfirm",
		Action:     "Click",
		Box:        maa.Rect{1, 2, 3, 4},
		Success:    true,
		DetailJson: `{"point":[2,3],"contact":0}`,
	}

	got := serializeActionDetail(detail)

	if got["id"] != int64(9) || got["name"] != "TapConfirm" || got["action"] != "Click" || got["success"] != true {
		t.Fatalf("unexpected action summary: %#v", got)
	}
	if detailJSON, ok := got["detail"].(map[string]interface{}); !ok || detailJSON["contact"] == nil {
		t.Fatalf("expected parsed detail JSON, got %#v", got["detail"])
	}
}
