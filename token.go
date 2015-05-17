package typo

// (c) Dmitriy Blokhin (sv.dblokhin@gmail.com), www.webjinn.ru

import (
    "io"
    "unicode/utf8"
)

type token struct {
    data string
    Is   lexType
}

func (l *token) len() int {
    return utf8.RuneCountInString(l.data)
}

func (l *token) Write(wr io.Writer) {
    wr.Write([]byte(l.data))
}