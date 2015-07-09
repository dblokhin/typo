package typo

// (c) Dmitriy Blokhin (sv.dblokhin@gmail.com), www.webjinn.ru

import (
    "strings"
)

func anSkipTag(an *analizer) fnState {

    current := an.scan()
    for current.Is != lexEof {
        current.Write(an.wr)
        if current.Is == lexCloseTag {
            return anOutput
        }

        current = an.scan()
    }

    // never reached
    return nil
}

func anOutput(an *analizer) fnState {

    current := an.scan()

    switch current.Is {
        case lexEof: {
            return nil
        }

        case lexSpace: {
            prev := an.get(-2)
            next := an.get(0)

            switch prev.Is {
                // Не отбивать пробел в начале текста
                case lexParagraph, lexEof: return anOutput;
                // Отбить неразрывный  пробел после коротких слов
                case lexWord: {
                    // если следом не пунктуация, конец
                    if prev.len() < 3 {
                        switch next.Is {
                            case lexWord, lexNumber, lexNo: {
                                // Исключения
                                switch strings.ToLower(prev.data) {
                                    case "же", "бы", "ли": {
                                        lspace.Write(an.wr)
                                        return anOutput
                                    }
                                }

                                // Если слово привязано к знаку пунктуации
                                prev := an.get(-3)
                                if prev.Is == lexPunctuation || prev.Is == lexHellip {
                                    lspace.Write(an.wr)
                                    return anOutput
                                }

                                lnbsp.Write(an.wr)
                                return anOutput
                            }
                        }
                    }
                }

                case lexNumber: {
                    // Отбить триады чисел thinsp
                    if next.Is == lexNumber && next.len() == 3 {
                        lthinsp.Write(an.wr)
                        return anOutput
                    }
                    // Не отбивать пробел перед %
                    if next.Is == lexPercent {
                        return anOutput
                    }
                }

                case lexSection: {
                    next := an.get(0)

                    // Отбивать неразрывным пробелом с текстом
                    if  next.Is == lexNumber || next.Is == lexWord {
                        lnbsp.Write(an.wr)
                    }

                    return anOutput
                }

            }

            switch next.Is {
                // Не отбивать пробел перед пунктуацией, концом текста, закрытой скобкой, тире
                case lexPunctuation, lexEof, lexParagraph, lexCloseBracket, lexDefisominus: return anOutput;
                case lexWord: {
                    switch strings.ToLower(next.data) {
                        case "же", "бы", "ли": {
                            lnbsp.Write(an.wr)
                            return anOutput
                        }
                    }
                }
            }

            current.Write(an.wr)
            return anOutput
        }

        case lexDefisominus: {
            prevns := an.getPrevNS(-2)
            prev := an.get(-2)
            next := an.get(0)

            if prev.Is != lexSpace && prev.Is != lexParagraph {
                // Если прямиком за словом, символом, точкой, числом и т.п. оставляем
                current.Write(an.wr)
            } else {
                switch prevns.Is {
                    // Выделение прямой речи после знаков препинания
                    case lexPunctuation, lexHellip: {
                        lspace.Write(an.wr)
                        lmdash.Write(an.wr)
                        lnbsp.Write(an.wr)  // TODO: Не печатать, если конец строки

                        // Пропустить следующий пробел, если есть
                        if next.Is == lexSpace {
                            an.peek()
                        }
                        return anOutput
                    }
                    // Выделение прямой речи в начале текста
                    case lexParagraph, lexEof: {
                        lmdash.Write(an.wr)
                        lnbsp.Write(an.wr)

                        // Пропустить следующий пробел, если есть
                        if next.Is == lexSpace {
                            an.peek()
                        }
                        return anOutput
                    }

                    // привязка к предыдущему слову
                    case lexWord, lexNumber, lexCloseTag, lexCloseBracket: {
                        lnbsp.Write(an.wr)
                        lmdash.Write(an.wr)
                        lspace.Write(an.wr)

                        // Пропустить следующий пробел, если есть
                        if next.Is == lexSpace {
                            an.peek()
                        }
                        return anOutput
                    }
                    // Во всех прочих случаях просто напечатать
		    default: {
		        current.Write(an.wr)
		        return anOutput
		    }
                
                }
            }

            return anOutput;
        }

        case lexPunctuation: {
            current.Write(an.wr)
            next := an.get(0)

            switch next.Is {
                case lexPunctuation, lexEof, lexParagraph: return anOutput;
            }

            return anOutput
        }

        case lexSection: {
            current.Write(an.wr)
            next := an.get(0)

            // Отбивать неразрывным пробелом с текстом
            if  next.Is == lexNumber || next.Is == lexWord {
                lnbsp.Write(an.wr)
            }

            return anOutput
        }

        case lexOpenBracket: {
            prev := an.get(-2)

            // Отбивать пробелом перед текстом
            if prev.Is == lexWord {
                lspace.Write(an.wr)
            }

            current.Write(an.wr)
            return anOutput
        }

        case lexCloseBracket: {
            current.Write(an.wr)

            // Для закрывающей скобки отбивать пробел после
            if current.Is == lexCloseBracket {
                next := an.get(0)
                switch next.Is {
                    case lexPunctuation, lexHellip, lexParagraph, lexEof, lexCloseBracket: return anOutput
                    default: lspace.Write(an.wr)
                }
            }

            return anOutput
        }

        case lexOpenTag: {
            current.Write(an.wr)
            return anSkipTag
        }

        case lexNo: {
            current.Write(an.wr)
            next := an.getNextNS(0)

            switch next.Is {
                case lexWord, lexNumber: {
                    lnbsp.Write(an.wr)

                    // Пропустить следующий пробел, если есть
                    next = an.get(0)
                    if next.Is == lexSpace {
                        an.peek()
                    }
                }
            }

            return anOutput
        }

        default: {
            current.Write(an.wr)
            return anOutput
        }
    }

    return nil
}

func anOpenQuote(an *analizer) fnState {
    // Сохраняю старую позицию
    current := an.scan()
    srartPos := an.pos

    // Определяю какая ковычка должна быть
    // Если за ней пробел, то quot
    if current.Is == lexQuote && current.data == "\"" {
        next := an.get(0)
        prev := an.get(-2)

        // Открывающая елочка
        if next.Is != lexSpace && next.Is != lexParagraph {
            an.update(-1, llaquo)
            an.seek(srartPos)
            return anPrepare
        }

        // Закрывающая елочка
        if prev.Is != lexSpace && prev.Is != lexParagraph {
            an.update(-1, lraquo)
            an.seek(srartPos)
            return anPrepare
        }

        an.update(-1, lquot)
        an.seek(srartPos)
        return anPrepare
    }

    return anPrepare
}

func anOpenAmpersand(an *analizer) fnState {

    next1 := an.get(1)
    if next1.Is == lexEof {
        an.peek()
        return anPrepare
    }
    next2 := an.get(2)
    if next2.Is == lexEof {
        an.peek()
        return anPrepare
    }

    if next2.data == ";" && next1.Is == lexWord {
        switch strings.ToLower(next1.data) {
            case "nbsp": {
                an.dropN(3); // Удалить неразрывный пробел
                return anPrepare
            }
            case "amp": {
                an.dropN(3);
                an.insert(lamp)
                an.peek()
                return anPrepare
            }
            case "mdash", "ndash": {
                an.dropN(3);
                an.insert(lexdefisominus)
                an.peek()
                return anPrepare
            }
            case "sect": {
                an.dropN(3);
                an.insert(lsect)
                an.peek()
                return anPrepare
            }
            case "quot", "raquo", "laquo": {
                an.dropN(3);
                an.insert(lexquote)
                an.peek()
                return anPrepare
            }
            case "hellip": {
                an.dropN(3);
                an.insert(lhellip)
                an.peek()
                return anPrepare
            }
            case "#8470": {
                an.dropN(3);
                an.insert(lno)
                an.peek()
                return anPrepare
            }

            default: an.update(0, lamp)
        }
    } else {
        an.update(0, lamp)
    }

    an.peek()
    return anPrepare
}

func anOpenBracket(an *analizer) fnState {
    next1 := an.get(1)
    if next1.Is == lexEof {
        an.peek()
        return anPrepare
    }
    next2 := an.get(2)
    if next2.Is == lexEof {
        an.peek()
        return anPrepare
    }

    if next2.data == ")" && next1.Is == lexWord {
        switch strings.ToLower(next1.data) {
            case "tm": {
                an.dropN(3); // Удалить неразрывный пробел
                an.insert(ltrade)
                an.peek()
                return anPrepare
            }
            case "r": {
                an.dropN(3);
                an.insert(lreg)
                an.peek()
                return anPrepare
            }
            case "c": {
                an.dropN(3);
                an.insert(lcopy)
                an.peek()
                return anPrepare
            }
            case "sect": {
                an.dropN(3);
                an.insert(lsect)
                an.peek()
                return anPrepare
            }

        }
    }

    an.peek()
    return anPrepare
}

func anCheckPunctuation(an *analizer) fnState {
    next0 := an.get(0)
    if next0.Is == lexEof {
        an.peek()
        return anPrepare
    }
    next1 := an.get(1)
    if next1.Is == lexEof {
        an.peek()
        return anPrepare
    }
    next2 := an.get(2)
    if next2.Is == lexEof {
        an.peek()
        return anPrepare
    }

    // Если многоточие:
    if next0.data == "." && next1.data == "." && next2.data == "." {
        an.dropN(3);
        an.insert(lhellip)
        an.peek()
        return anPrepare
    }

    an.peek()
    return anPrepare
}

func anStartSpace(an *analizer) fnState {
    next := an.get(1)
    prev := an.get(-1)
    // Удалить повторяющиеся пробелы и те, что в начале абзаца
    if next.Is == lexSpace  || prev.Is == lexParagraph {
        an.drop(0)
        return anPrepare
    }

    an.peek()
    return anPrepare
}

// anStart подготовка набора лексем (упаковка, опеределение чисел, сокращений, знаков т.п.)
func anPrepare(an *analizer) fnState {

    current := an.get(0)

    switch current.Is {
        case lexEof: {
            return nil
        }

        case lexSpace: {
            return anStartSpace
        }

        case lexPunctuation: {
            return anCheckPunctuation
        }

        case lexAmpersand: {
            return anOpenAmpersand
        }

        case lexOpenBracket: {
            return anOpenBracket
        }

        case lexQuote: {
            if current.data != "\"" {
                an.peek()
                return anPrepare
            }

            return anOpenQuote
        }

        case lexWord: {
            if isDecimal(current.data) {
                current.Is = lexNumber
                return anPrepare
            }

            an.peek()
            return anPrepare
        }

        default: {
            an.peek()
            return anPrepare
        }
    }

    return nil
}

