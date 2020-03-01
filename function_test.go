package dakuten

import "testing"

func TestDakuon(t *testing.T) {
	input := "たねだりさ。種田梨沙。Taneda Risa. ﾀﾈﾀﾞﾘｻ。ポピパ、ぽぴぱ、ﾎﾟﾋﾟﾊﾟ!"
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

// はひふへほハヒフヘホ
// ぱぴぷぺぽパピプペポ
// 〳〱〴〲
func TestDakuon_Mapping(t *testing.T) {
	input := "かきくけこさしすせそたちつてとはひふへほうゝカキクケコサシスセソタチツテトハヒフヘホウワヰヱヲヽがぎぐげござじずぜぞだぢづでどばびぶべぼゔゞガギグゲゴザジズゼゾダヂヅデドバビブベボヴヷヸヹヺヾぱぴぷぺぽパピプペポ"

	dakuonExpected := "がぎぐげござじずぜぞだぢづでどばびぶべぼゔゞガギグゲゴザジズゼゾダヂヅデドバビブベボヴヷヸヹヺヾがぎぐげござじずぜぞだぢづでどばびぶべぼゔゞガギグゲゴザジズゼゾダヂヅデドバビブベボヴヷヸヹヺヾばびぶべぼバビブベボ"
	dakuonGot := toCombiningDakuon(input)
	if dakuonGot != dakuonExpected {
		t.Errorf("toCombiningDakuon; want %s; got %s", dakuonExpected, dakuonGot)
	}

	handakuonExpected := "か゚き゚く゚け゚こ゚さ゚し゚す゚せ゚そ゚た゚ち゚つ゚て゚と゚ぱぴぷぺぽう゚ゝ゚カ゚キ゚ク゚ケ゚コ゚サ゚シ゚ス゚セ゚ソ゚タ゚チ゚ツ゚テ゚ト゚パピプペポウ゚ワ゚ヰ゚ヱ゚ヲ゚ヽ゚か゚き゚く゚け゚こ゚さ゚し゚す゚せ゚そ゚た゚ち゚つ゚て゚と゚ぱぴぷぺぽう゚ゝ゚カ゚キ゚ク゚ケ゚コ゚サ゚シ゚ス゚セ゚ソ゚タ゚チ゚ツ゚テ゚ト゚パピプペポウ゚ワ゚ヰ゚ヱ゚ヲ゚ヽ゚ぱぴぷぺぽパピプペポ"
	handakuonGot := toCombiningHandakuon(input)
	if handakuonGot != handakuonExpected {
		t.Errorf("toCombiningHandakuon; want %s; got %s", handakuonExpected, handakuonGot)
	}
}

func TestDakuon_Odoriji(t *testing.T) {
	input := "〳〱〴〲"

	dakuonExpected := "〴〲〴〲"
	dakuonGot := toCombiningDakuon(input)
	if dakuonGot != dakuonExpected {
		t.Errorf("toCombiningDakuon; want %s; got %s", dakuonExpected, dakuonGot)
	}

	handakuonExpected := "〳゚〱゚〳゚〱゚"
	handakuonGot := toCombiningHandakuon(input)
	if handakuonGot != handakuonExpected {
		t.Errorf("toCombiningHandakuon; want %s; got %s", handakuonExpected, handakuonGot)
	}
}

func TestDakuon_Control(t *testing.T) {
	input := "\n\t いちえちゃん"
	dakuonExpected := "\n\t ﾞい゙ぢえ゙ぢゃ゙ん゙"
	dakuonGot := toCombiningDakuon(input)
	if dakuonGot != dakuonExpected {
		t.Errorf("toCombiningDakuon; want %s; got %s", dakuonExpected, dakuonGot)
	}
}

func TestDakuon_Combining(t *testing.T) {
	input := "🇯🇵"
	dakuonExpected := "🇯🇵ﾞ"
	dakuonGot := toCombiningDakuon(input)
	if dakuonGot != dakuonExpected {
		t.Errorf("toCombiningDakuon; want %s; got %s", dakuonExpected, dakuonGot)
	}
}
