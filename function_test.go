package dakuten

import "testing"

const input = "たねだりさ。種田梨沙。Taneda Risa. ﾀﾈﾀﾞﾘｻ。ポピパ、ぽぴぱ、ﾎﾟﾋﾟﾊﾟ!"

func TestDakuon(t *testing.T) {
	tests := []struct {
		name     string
		f        (func(string) string)
		expected string
	}{
		{"toExternalDakuon", toExternalDakuon, "た゛ね゛た゛り゛さ゛。゛種゛田゛梨゛沙゛。゛Tﾞaﾞnﾞeﾞdﾞaﾞ ﾞRﾞiﾞsﾞaﾞ.ﾞ ﾞﾀﾞﾈﾞﾀﾞﾘﾞｻﾞ。゛ホ゛ヒ゛ハ゛、゛ほ゛ひ゛は゛、゛ﾎﾞﾋﾞﾊﾞ!ﾞ"},
		{"toCombiningDakuon", toCombiningDakuon, "だね゙だり゙ざ。゙種゙田゙梨゙沙゙。゙Tﾞaﾞnﾞeﾞdﾞaﾞ ﾞRﾞiﾞsﾞaﾞ.ﾞ ﾞﾀﾞﾈﾞﾀﾞﾘﾞｻﾞ。゙ボビバ、゙ぼびば、゙ﾎﾞﾋﾞﾊﾞ!ﾞ"},
		{"toExternalHandakuon", toExternalHandakuon, "た゜ね゜た゜り゜さ゜。゜種゜田゜梨゜沙゜。゜Tﾟaﾟnﾟeﾟdﾟaﾟ ﾟRﾟiﾟsﾟaﾟ.ﾟ ﾟﾀﾟﾈﾟﾀﾟﾘﾟｻﾟ。゜ホ゜ヒ゜ハ゜、゜ほ゜ひ゜は゜、゜ﾎﾟﾋﾟﾊﾟ!ﾟ"},
		{"toCombiningHandakuon", toCombiningHandakuon, "た゚ね゚た゚り゚さ゚。゚種゚田゚梨゚沙゚。゚Tﾟaﾟnﾟeﾟdﾟaﾟ ﾟRﾟiﾟsﾟaﾟ.ﾟ ﾟﾀﾟﾈﾟﾀﾟﾘﾟｻﾟ。゚ポピパ、゚ぽぴぱ、゚ﾎﾟﾋﾟﾊﾟ!ﾟ"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.f(input)
			if got != tt.expected {
				t.Errorf("%s; want %s; got %s", tt.name, tt.expected, got)
			}
		})
	}
}
