package pdf_export

import "testing"

func TestLabelsFor_Romanian(t *testing.T) {
	l := labelsFor("ro")
	if l.Food != "Aliment" || l.Quantity != "Cantitate" || l.Alternative != "Alternativă" || l.Supplements != "Suplimente" {
		t.Errorf("unexpected Romanian labels: %+v", l)
	}
}

func TestLabelsFor_English(t *testing.T) {
	l := labelsFor("en")
	if l.Food != "Food" || l.Quantity != "Quantity" || l.Alternative != "Alternative" || l.Supplements != "Supplements" {
		t.Errorf("unexpected English labels: %+v", l)
	}
}

func TestLabelsFor_EmptyDefaultsToRomanian(t *testing.T) {
	l := labelsFor("")
	if l != labelsFor("ro") {
		t.Errorf("expected empty language to default to Romanian, got %+v", l)
	}
}

func TestLabelsFor_UnknownDefaultsToRomanian(t *testing.T) {
	l := labelsFor("fr")
	if l != labelsFor("ro") {
		t.Errorf("expected an unrecognized language to default to Romanian, got %+v", l)
	}
}
