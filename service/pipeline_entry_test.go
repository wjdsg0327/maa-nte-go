package service

import "testing"

func TestResolvePipelineEntryUsesSameNamedNode(t *testing.T) {
	nodes := map[string]interface{}{
		"F1": map[string]interface{}{
			"recognition": "DirectHit",
		},
		"Other": map[string]interface{}{
			"recognition": "DirectHit",
		},
	}

	entry, err := resolvePipelineEntry("F1", nodes)
	if err != nil {
		t.Fatalf("resolve entry: %v", err)
	}
	if entry != "F1" {
		t.Fatalf("expected same-named entry, got %q", entry)
	}
}

func TestResolvePipelineEntryUsesUniqueRootWhenFileNameIsNotNodeName(t *testing.T) {
	nodes := map[string]interface{}{
		"点击F1": map[string]interface{}{
			"recognition": "DirectHit",
			"next":        "点击活跃度",
		},
		"点击活跃度": map[string]interface{}{
			"recognition": "TemplateMatch",
		},
	}

	entry, err := resolvePipelineEntry("F1", nodes)
	if err != nil {
		t.Fatalf("resolve entry: %v", err)
	}
	if entry != "点击F1" {
		t.Fatalf("expected root entry 点击F1, got %q", entry)
	}
}

func TestResolvePipelineEntryRejectsAmbiguousRoots(t *testing.T) {
	nodes := map[string]interface{}{
		"A": map[string]interface{}{"recognition": "DirectHit"},
		"B": map[string]interface{}{"recognition": "DirectHit"},
	}

	_, err := resolvePipelineEntry("F1", nodes)
	if err == nil {
		t.Fatal("expected ambiguous roots to return an error")
	}
}

func TestRelationTargetNamesSupportStringListAndNodeAttr(t *testing.T) {
	value := []interface{}{
		"A",
		map[string]interface{}{"name": "B", "jump_back": true},
		map[string]interface{}{"name": "Anchor", "anchor": true},
	}

	got := relationTargetNames(value)
	if len(got) != 2 || got[0] != "A" || got[1] != "B" {
		t.Fatalf("expected concrete targets A and B, got %#v", got)
	}
}
