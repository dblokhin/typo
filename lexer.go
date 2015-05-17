package typo

// (c) Dmitriy Blokhin (sv.dblokhin@gmail.com), www.webjinn.ru

import (
    "errors"
    "unicode/utf8"
)

var (
    terminals = map[rune]bool{
        '.': true, ',': true, '?': true, '!': true, ';': true, '%': true,
        ':': true, '-': true, '\'': true, '"': true, ' ': true, '\t': true,
        '(': true, ')': true, '&': true, '\n': true, '№': true, '§': true,
        '<': true, '>': true, '™': true, '©': true, '®': true,
        '…': true,
    }

    tokens = map[rune]bool{
        '.': true, ',': true, '?': true, '!': true, ';': true, '%': true,
        ':': true, '-': true, '\'': true, '"': true, ' ': true, '\t': true,
        '(': true, ')': true, '&': true, '\n': true, '№': true, '§': true,
        '<': true, '>': true, '™': true, '©': true, '®': true,
        '…': true,
    }

    ErrEofParser = errors.New("EOF parser")
    ErrRuneError = errors.New("Rune error")

    lexdefisominus = &token{"-", lexDefisominus}
    lexquote = &token{"\"", lexQuote}

    lspace = &token{" ", lexSpace}
    lnbsp = &token{"&nbsp;", lexSpace}
    ldot = &token{".", lexPunctuation}
    lrsquo = &token{"&rsquo;", lexWord}
    lquot = &token{"&quot;", lexQuote}
    lraquo = &token{"&raquo;", lexQuote}
    llaquo = &token{"&laquo;", lexQuote}
    lbdquo = &token{"&bdquo;", lexQuote}
    lldquo = &token{"&ldquo;", lexQuote}
    lthinsp = &token{"&thinsp;", lexDefisominus}
    lmdash = &token{"&mdash;", lexDefisominus}
    lndash = &token{"&ndash;", lexDefisominus}
    lamp = &token{"&amp;", lexWord}
    lsect = &token{"&sect;", lexSection}
    lhellip = &token{"&hellip;", lexHellip}
    lcopy = &token{"&copy;", lexWord}
    ltrade = &token{"&trade;", lexWord}
    lreg = &token{"&reg;", lexWord}
    lno = &token{"&#8470;", lexNo}
    lparagraph  = &token{"", lexParagraph}
    leof = &token{"", lexEof}
)

type lexType int

const (
    lexSpace lexType = iota
    lexParagraph
    lexPunctuation                 // '.', ',', '?', '!', ';', ':', '...'
    lexWord
    lexNumber
    lexDefisominus
    lexPercent
    lexSection
    lexHellip
    lexNo

    lexQuote                      // "
    lexOpenBracket                  // (
    lexCloseBracket                 // )
    lexOpenTag                  // <
    lexCloseTag                 // >
    lexAmpersand                // &
    lexCelsi                    // C



    lexEof
)

type lexer struct {
    data string
    pos  int
    tokens []*token
}

// newLexer возвращает Lexer с разобранным набором лексем
func newLexer(str string) *lexer {
    tokens := make([]*token, 0)
    lex := &lexer{str, 0, tokens}
    lex.run()

    return lex
}

// run Сканирует исходный текст и собирает массив лексем
func (l *lexer) run() {
    // Добавляет код начала параграфа
    l.addToken(lparagraph)

    // Обработка текста
    for {
        token, err := l.nextToken()

        if err != nil {
            break
        }

        l.addToken(token)
    }
}

// nextToken возвращает следующий токен (лексему)
func (l *lexer) nextToken() (result *token, err error) {

    char, err := l.scan()
    if err != nil {
        return nil, err
    }

    // Пропускаем пробелы
    if char == ' ' || char == '\t' || char == '\r' {
        return lspace, nil
    }

    // Если односимволный токен
    if isToken(char) {
        return tokenByRune(char), nil
    }

    // Читаем токен до следующего терминала
    str := ""
    for {
        str = str + string(char)

        char, err = l.next()
        if isTerminal(char) || err == ErrEofParser {
            break
        }

        if err != nil {
            return nil, err
        }

        l.pick()
    }

    return &token{str, lexWord}, nil
}

// scan читает Rune из входного потока и возвращает его, Pos обновляет
func (l *lexer) scan() (rune, error) {
    if l.pos >= len(l.data) {
        return 0, ErrEofParser
    }

    r, s := utf8.DecodeRune([]byte(l.data[l.pos:]))
    if r == utf8.RuneError {
        return 0, ErrRuneError
    }
    l.pos += s

    return r, nil
}

// pick перемещает позицию на 1 Rune вперед
func (l *lexer) pick() {
    l.scan()
}

// next читает Rune из входного потока и возвращает его, Pos не обновляет
func (l *lexer) next() (rune, error) {
    if l.pos >= len(l.data) {
        return 0, ErrEofParser
    }

    r, _ := utf8.DecodeRune([]byte(l.data[l.pos:]))
    if r == utf8.RuneError {
        return 0, ErrRuneError
    }

    return r, nil
}

// addToken добавляет токен(лексему) в массив токенов
func (l *lexer) addToken(t *token) {
    l.tokens = append(l.tokens, t)
}

// isTerminal возвращает true если символ терминальный
func isTerminal(val rune) bool {
    _, ok := terminals[val]

    return ok
}

// isToken возвращает true если символ - токен
func isToken(val rune) bool {
    _, ok := tokens[val]

    return ok
}

// tokenByRune возвращает токен по символу (если это символ-токен)
func tokenByRune(val rune) *token {
    switch val {
        case '.', ',', '?', '!', ';', ':':
        {
            return &token{string(val), lexPunctuation}
        }
        case '\'':
        {
            return lrsquo
        }
        case '"', '«', '»':
        {
            return &token{"\"", lexQuote}
        }
        case '%':
        {
            return &token{string(val), lexPercent}
        }
        case '(': {
            return &token{string(val), lexOpenBracket}
        }
        case ')': {
            return &token{string(val), lexCloseBracket}
        }
        case '<': {
            return &token{string(val), lexOpenTag}
        }
        case '>': {
            return &token{string(val), lexCloseTag}
        }
        case '&': {
            return &token{string(val), lexAmpersand}
        }
        case '-': {
            return &token{string(val), lexDefisominus}
        }
        case '§': {
            return lsect
        }
        case '…': {
            return lhellip
        }
        case '®': {
            return lreg
        }
        case '©': {
            return lcopy
        }
        case '™': {
            return ltrade
        }
        case '№': {
            return lno
        }
        case '\n': {
            return lparagraph
        }
    }

    return &token{"", lexEof}
}