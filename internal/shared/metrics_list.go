package shared

import (
	"github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

type SliceMetrics []*Metrics

func (ml SliceMetrics) MarshalJSON() ([]byte, error) {
	w := &jwriter.Writer{}
	customMarshalEasyJSON(w, &ml)
	return w.Buffer.BuildBytes(), nil
}

func (ml *SliceMetrics) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	customUnmarshalEasyJSON(&r, ml)
	return r.Error()
}

func customMarshalEasyJSON(w *jwriter.Writer, ml *SliceMetrics) {
	w.RawByte('[')
	for i, metric := range *ml {
		if i > 0 {
			w.RawByte(',')
		}
		w.Raw(easyjson.Marshal(metric))
	}
	w.RawByte(']')
}

func customUnmarshalEasyJSON(in *jlexer.Lexer, out *SliceMetrics) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('[')
	for !in.IsDelim(']') {
		var v1 *Metrics
		if in.IsNull() {
			in.Skip()
			v1 = nil
		} else {
			v1 = new(Metrics)
			(*v1).UnmarshalEasyJSON(in)
		}
		*out = append(*out, v1)
		in.WantComma()
	}
	in.Delim(']')
	in.WantComma()
	if isTopLevel {
		in.Consumed()
	}
}

func (ml SliceMetrics) MarshalEasyJSON(w *jwriter.Writer) {
	customMarshalEasyJSON(w, &ml)
}

func (ml *SliceMetrics) UnmarshalEasyJSON(l *jlexer.Lexer) {
	customUnmarshalEasyJSON(l, ml)
}
