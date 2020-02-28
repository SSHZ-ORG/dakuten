package dakuten

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/text/width"
)

var (
	dakuonTable = map[rune]rune{
		'か': 'が', 'き': 'ぎ', 'く': 'ぐ', 'け': 'げ', 'こ': 'ご',
		'さ': 'ざ', 'し': 'じ', 'す': 'ず', 'せ': 'ぜ', 'そ': 'ぞ',
		'た': 'だ', 'ち': 'ぢ', 'つ': 'づ', 'て': 'で', 'と': 'ど',
		'は': 'ば', 'ひ': 'び', 'ふ': 'ぶ', 'へ': 'べ', 'ほ': 'ぼ',
		'う': 'ゔ',
		'ゝ': 'ゞ',
		'カ': 'ガ', 'キ': 'ギ', 'ク': 'グ', 'ケ': 'ゲ', 'コ': 'ゴ',
		'サ': 'ザ', 'シ': 'ジ', 'ス': 'ズ', 'セ': 'ゼ', 'ソ': 'ゾ',
		'タ': 'ダ', 'チ': 'ヂ', 'ツ': 'ヅ', 'テ': 'デ', 'ト': 'ド',
		'ハ': 'バ', 'ヒ': 'ビ', 'フ': 'ブ', 'ヘ': 'ベ', 'ホ': 'ボ',
		'ウ': 'ヴ',
		'ワ': 'ヷ', 'ヰ': 'ヸ', 'ヱ': 'ヹ', 'ヲ': 'ヺ',
		'ヽ': 'ヾ',
		'〳': '〴', '〱': '〲',
	}
	reversedDakuonTable = reversed(dakuonTable)

	handakuonTable = map[rune]rune{
		'は': 'ぱ', 'ひ': 'ぴ', 'ふ': 'ぷ', 'へ': 'ぺ', 'ほ': 'ぽ',
		'ハ': 'パ', 'ヒ': 'ピ', 'フ': 'プ', 'ヘ': 'ペ', 'ホ': 'ポ',
	}
	reversedHandakuonTable = reversed(handakuonTable)
)

func reversed(m map[rune]rune) map[rune]rune {
	o := make(map[rune]rune)
	for k, v := range m {
		o[v] = k
	}
	return o
}

const (
	fdm = '゛' // U+309B, KATAKANA-HIRAGANA VOICED SOUND MARK
	cdm = '゙' // U+3099, COMBINING KATAKANA-HIRAGANA VOICED SOUND MARK
	hdm = 'ﾞ' // U+FF9E, HALFWIDTH KATAKANA VOICED SOUND MARK

	fhm = '゜' // U+309C, KATAKANA-HIRAGANA SEMI-VOICED SOUND MARK
	chm = '゚' // U+309A, COMBINING KATAKANA-HIRAGANA SEMI-VOICED SOUND MARK
	hhm = 'ﾟ' // U+FF9F, HALFWIDTH KATAKANA SEMI-VOICED SOUND MARK
)

func toExternalInternal(in string, fm, hm rune) string {
	o := strings.Builder{}
	for _, i := range in {
		if i == fdm || i == cdm || i == hdm || i == fhm || i == chm || i == hhm {
			continue
		}

		if e, ok := reversedDakuonTable[i]; ok {
			i = e
		} else if e, ok := reversedHandakuonTable[i]; ok {
			i = e
		}

		o.WriteRune(i)
		switch width.LookupRune(i).Kind() {
		case width.EastAsianFullwidth, width.EastAsianWide:
			o.WriteRune(fm)
		default:
			o.WriteRune(hm)
		}
	}
	return o.String()
}

func toExternalDakuon(in string) string {
	return toExternalInternal(in, fdm, hdm)
}

func toExternalHandakuon(in string) string {
	return toExternalInternal(in, fhm, hhm)
}

func toCombiningInternal(in string, t map[rune]rune, cm, hm rune) string {
	o := strings.Builder{}
	for _, i := range in {
		if i == fdm || i == cdm || i == hdm || i == fhm || i == chm || i == hhm {
			continue
		}

		if e, ok := reversedDakuonTable[i]; ok {
			i = e
		} else if e, ok := reversedHandakuonTable[i]; ok {
			i = e
		}

		if e, ok := t[i]; ok {
			o.WriteRune(e)
			continue
		}

		o.WriteRune(i)
		switch width.LookupRune(i).Kind() {
		case width.EastAsianFullwidth, width.EastAsianWide:
			o.WriteRune(cm)
		default:
			o.WriteRune(hm)
		}
	}
	return o.String()
}

func toCombiningDakuon(in string) string {
	return toCombiningInternal(in, dakuonTable, cdm, hdm)
}

func toCombiningHandakuon(in string) string {
	return toCombiningInternal(in, handakuonTable, chm, hhm)
}

type answerInlineQueryResponse struct {
	Method        string                              `json:"method"` // Must be answerInlineQuery
	InlineQueryID string                              `json:"inline_query_id"`
	Results       []tgbotapi.InlineQueryResultArticle `json:"results"`
}

func newInlineQueryResultArticle(id, title, content string) tgbotapi.InlineQueryResultArticle {
	r := tgbotapi.NewInlineQueryResultArticle(id, title, content)
	if len(content) <= 64 {
		r.Description = content
	} else {
		last := 0
		for i := range content {
			if i <= 64 {
				last = i
			} else {
				break
			}
		}
		r.Description = content[:last]
	}
	return r
}

func Webhook(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)

	var update tgbotapi.Update
	err := json.Unmarshal(bytes, &update)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if update.InlineQuery == nil || update.InlineQuery.Query == "" {
		return
	}

	id := update.InlineQuery.ID
	query := update.InlineQuery.Query

	result := answerInlineQueryResponse{
		Method:        "answerInlineQuery",
		InlineQueryID: id,
		Results: []tgbotapi.InlineQueryResultArticle{
			newInlineQueryResultArticle(id+"dc", "濁音（結合文字）", toCombiningDakuon(query)),
			newInlineQueryResultArticle(id+"hc", "半濁音（結合文字）", toCombiningHandakuon(query)),
			newInlineQueryResultArticle(id+"de", "濁音", toExternalDakuon(query)),
			newInlineQueryResultArticle(id+"he", "半濁音", toExternalHandakuon(query)),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
