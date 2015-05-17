package typo

import (
    "io"
)

// (c) Dmitriy Blokhin (sv.dblokhin@gmail.com), www.webjinn.ru

type fnState func(*analizer) fnState

// analizer набора лексем и методов их анализа
type analizer struct {
    tokens []*token
    count, pos int
    wr io.Writer
}

func newAnalizer(tokens []*token) *analizer {
    var result analizer

    result.tokens = tokens
    result.pos = 0
    result.count = len(tokens)

    return &result
}

// run запускает анализ лексем
func (an *analizer) run(fn fnState) {
    an.pos = 0

    for state := fn; state != nil; {
        state = state(an)
    }
}

// update изменяет idx-ый токен относительно текущей позиции (pos) на t
func (an *analizer) update(idx int, t *token) {
    k := an.pos + idx

    if (k < an.count && k >= 0) {
        an.tokens[k] = t
    }
}

// getNextNoSpace возвращает следующий токен относительно текущей позиции (pos), пропуская пробелы
func (an *analizer) getNextNS(idx int) (*token) {
    t := an.get(idx)

    for t.Is == lexSpace {
        idx++
        t = an.get(idx)
    }

    return t
}

// getPrevNoSpace возвращает предыдущий токен относительно текущей позиции (pos), пропуская пробелы
func (an *analizer) getPrevNS(idx int) (*token) {
    t := an.get(idx)

    for t.Is == lexSpace {
        idx--
        t = an.get(idx)
    }

    return t
}

// get возвращает idx-ый токен относительно текущей позиции (pos)
func (an *analizer) get(idx int) (*token) {
    k := an.pos + idx

    if (k < an.count && k >= 0) {
        return an.tokens[k]
    }

    return leof
}

// scan читает из входного потока токен и перемещает позицию (pos)
func (an *analizer) scan() (*token) {
    token := an.get(0)
    if token.Is != lexEof {
        an.peek()
    }

    return token
}

// peek увеличивает позицию
func (an *analizer) peek() {
    an.pos++
}

// seek устанавливает позицию
func (an *analizer) seek(pos int) {
    an.pos = pos
}

// peekN увеличивает позицию на N
func (an *analizer) peekN(n int) {
    an.pos += n
}

// drop удаляет idx-ый токен относительно текущей позиции (pos)
func (an *analizer) drop(idx int) {
    k := an.pos + idx

    if (k < an.count && k >= 0) {
        an.tokens = append(an.tokens[:k], an.tokens[k + 1:]...)
    }

    an.count = len(an.tokens)
}

// drop удаляет n токенов относительно текущей позиции (pos)
func (an *analizer) dropN(n int) {
    for n > 0 {
        an.drop(0)
        n--
    }

    an.count = len(an.tokens)
}

// insert вставляет токен в текущую позицию (pos)
func (an *analizer) insert(t *token) {
    a := an.tokens
    an.tokens = append(a[:an.pos], append([]*token{t}, a[an.pos:]...)...)
    an.count = len(an.tokens)
}

func isDecimal(str string) bool {
    for i := 0; i < len(str); i++ {
        if '0' > str[i] || str[i] > '9' {
            return false
        }
    }

    return true
}