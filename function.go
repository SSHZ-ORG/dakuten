package dakuten

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

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

		if !unicode.IsGraphic(i) {
			o.WriteRune(i)
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

		if !unicode.IsGraphic(i) {
			o.WriteRune(i)
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

var converters = []struct {
	ID, Name string
	Func     func(string) string
}{
	{"dc", "濁点（結合文字）", toCombiningDakuon},
	{"hc", "半濁点（結合文字）", toCombiningHandakuon},
	{"de", "濁点", toExternalDakuon},
	{"he", "半濁点", toExternalHandakuon},
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

func handleInlineQuery(in *tgbotapi.InlineQuery) *answerInlineQueryResponse {
	id := in.ID
	query := in.Query

	results := []tgbotapi.InlineQueryResultArticle{}
	if query != "" {
		for _, c := range converters {
			results = append(results, newInlineQueryResultArticle(id+c.ID, c.Name, c.Func(query)))
		}
	}

	return &answerInlineQueryResponse{
		Method:        "answerInlineQuery",
		InlineQueryID: id,
		Results:       results,
	}
}

type sendMessageResponse struct {
	Method string `json:"method"` // Must be sendMessage
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func handleMessage(in *tgbotapi.Message) *sendMessageResponse {
	if in.Command() != "" {
		return nil
	}

	msg := in.Text
	if msg == "" {
		return nil
	}

	out := strings.Builder{}
	curLen := 0

	for _, c := range converters {
		s := c.Func(msg)
		thisLen := 1 + utf8.RuneCountInString(c.Name) + 1 + utf8.RuneCountInString(s) + 1
		if (curLen + thisLen) < 4096 {
			curLen += thisLen
			out.WriteRune('\n')
			out.WriteString(c.Name)
			out.WriteRune('\n')
			out.WriteString(s)
			out.WriteRune('\n')
		} else {
			break
		}
	}

	r := "Input is too long!"
	if out.Len() > 0 {
		r = out.String()
	}

	return &sendMessageResponse{
		Method: "sendMessage",
		ChatID: in.Chat.ID,
		Text:   r,
	}
}

func Webhook(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)

	var update tgbotapi.Update
	err := json.Unmarshal(bytes, &update)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var resp interface{} = nil
	if update.InlineQuery != nil {
		resp = handleInlineQuery(update.InlineQuery)
	} else if update.Message != nil {
		resp = handleMessage(update.Message)
	}

	if resp != nil {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}
