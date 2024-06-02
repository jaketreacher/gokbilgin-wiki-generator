package turkishsuffix

import (
	"testing"
)

func TestAblative(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Nejat Göyünç", "Nejat Göyünç'ten"},
		{"Halil Inalcik", "Halil Inalcik'ten"},
		{"Helmuth Scheel", "Helmuth Scheel'den"},
		{"Berberian Hrant", "Berberian Hrant'tan"},
		{"Pertev Boratav", "Pertev Boratav'dan"},
		{"Abramowicz Zygmunt", "Abramowicz Zygmunt'tan"},
		{"Sāṭi` al-Ḥuṣrī", "Sāṭi` al-Ḥuṣrī'den"},
	}

	for _, tt := range tests {
		result, err := Ablative(tt.input)
		if err != nil {
			t.Fatalf("Ablative(\"%s\") = \"%s\"; want \"%s\"", tt.input, err.Error(), tt.expected)
		}
		if result != tt.expected {
			t.Fatalf("Ablative(\"%s\") = \"%s\"; want \"%s\"", tt.input, result, tt.expected)
		}
	}
}
