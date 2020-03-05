package dakuten

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rivo/uniseg"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/width"
)

var (
	// Replacers to handle runes that are not properly handled by norm.NFD / norm.NFC.
	toNFCExtraReplacer, toNFDExtraReplacer = constructReplacers(map[string]string{
		// Odoriji
		"〳゙": "〴", // U+3033 + U+3099 = U+3034, VERTICAL KANA REPEAT MARK UPPER HALF
		"〱゙": "〲", // U+3031 + U+3099 = U+3032, VERTICAL KANA REPEAT MARK
	})
)

func constructReplacers(m map[string]string) (*strings.Replacer, *strings.Replacer) {
	var toC, toD []string
	for d, c := range m {
		toC = append(toC, d, c)
		toD = append(toD, c, d)
	}
	return strings.NewReplacer(toC...), strings.NewReplacer(toD...)
}

const (
	fdm = '゛' // U+309B, KATAKANA-HIRAGANA VOICED SOUND MARK
	cdm = '゙' // U+3099, COMBINING KATAKANA-HIRAGANA VOICED SOUND MARK
	hdm = 'ﾞ' // U+FF9E, HALFWIDTH KATAKANA VOICED SOUND MARK

	fhm = '゜' // U+309C, KATAKANA-HIRAGANA SEMI-VOICED SOUND MARK
	chm = '゚' // U+309A, COMBINING KATAKANA-HIRAGANA SEMI-VOICED SOUND MARK
	hhm = 'ﾟ' // U+FF9F, HALFWIDTH KATAKANA SEMI-VOICED SOUND MARK
)

func convertInternal(in string, fm, hm rune) string {
	o := strings.Builder{}

	gr := uniseg.NewGraphemes(toNFDExtraReplacer.Replace(norm.NFD.String(in)))
	for gr.Next() {
		rs := gr.Runes()
		r := rs[0]

		// Not printable. We should keep it as-is.
		if !unicode.IsGraphic(r) {
			o.WriteString(string(rs))
			continue
		}

		for _, i := range rs {
			if i == fdm || i == cdm || i == hdm || i == fhm || i == chm || i == hhm {
				continue
			}
			o.WriteRune(i)
		}

		switch width.LookupRune(r).Kind() {
		case width.EastAsianFullwidth, width.EastAsianWide:
			o.WriteRune(fm)
		default:
			o.WriteRune(hm)
		}
	}

	return toNFCExtraReplacer.Replace(norm.NFC.String(o.String()))
}

func toExternalDakuon(in string) string {
	return convertInternal(in, fdm, hdm)
}

func toExternalHandakuon(in string) string {
	return convertInternal(in, fhm, hhm)
}

func toCombiningDakuon(in string) string {
	return convertInternal(in, cdm, hdm)
}

func toCombiningHandakuon(in string) string {
	return convertInternal(in, chm, hhm)
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
