package typo

// (c) Dmitriy Blokhin (sv.dblokhin@gmail.com), www.webjinn.ru

import (
    "fmt"
    "log"
    "bytes"
    "io"
)

type Typograph struct {
    input string
    wr io.Writer
}

func New(str string, writer io.Writer) *Typograph {
    return &Typograph{str, writer}
}

func (t *Typograph) Typo() {
    // Подготовить массив лексем
    lex := newLexer(t.input)
    lex.run()

    // Первичная обработка массива лексем
    an := newAnalizer(lex.tokens)
    an.run(anPrepare)

    // Вывод
    an.wr = t.wr
    an.run(anOutput)
}

func Typo(str string) string {
    // Подготовить массив лексем
    lex := newLexer(str)
    lex.run()

    // Первичная обработка массива лексем
    an := newAnalizer(lex.tokens)
    an.run(anPrepare)

    // Вывод
    buf := new(bytes.Buffer)
    an.wr = buf
    an.run(anOutput)

    return buf.String()
}


















































/*


    for idx < len(pars.tokens) - 2 {
        current := pars.tokens[idx]
        next1 := pars.tokens[idx + 1]
        next2 := pars.tokens[idx + 2]

        reuse:
        switch current.Is {
            case lexDefisominus: {
                if (next1.Is == lexSpace && next2.Is == lexWord) {  // Выделение прямой речи
                    lmdash.Write(wr)
                    lnbsp.Write(wr)
                    idx++
                }
            }
            }
            case lexQuote: {
                if next1.Is == lexSpace {
                    lquot.Write(wr)
                } else {}
            }

}
*/


func delme() {
    fmt.Println("Ya!")
    log.Println("Ya!")
}