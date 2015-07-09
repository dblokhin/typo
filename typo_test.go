package typo

import (
    "testing"
)

// (c) Dmitriy Blokhin (sv.dblokhin@gmail.com), www.webjinn.ru

type InOUtStr struct {
    In, Out string
}

var (
    testerTypo []InOUtStr = []InOUtStr{
        {"", ""},
        {" ", ""},
        {"   ", ""},
        {"\t\t", ""},
        {"\t ", ""},
        {" text", "text"},
        {"text ", "text"},
        {" text ", "text"},
        {"   text", "text"},
        {"   text   ", "text"},
        {"text text text", "text text text"},
        {"text   text", "text text"},
        {"text   text   text", "text text text"},
        {"text \t\t text  \t", "text text"},
        {".text", ".text"},
        {". text", ". text"},
        {".text . text", ".text. text"},
        {"text... text", "text&hellip; text"},
        {"text?! text", "text?! text"},
        {"35 %", "35%"},{"35   %", "35%"},
        {"333 %", "333%"},
        {"22 22", "22 22"},
        {"2 333 22", "2&thinsp;333 22"},
        {"a222 333 222", "a222 333&thinsp;222"},
        {"- Я пошёл домой... - Может останешься? - Нет, ухожу.", "&mdash;&nbsp;Я&nbsp;пошёл домой&hellip; &mdash;&nbsp;Может останешься? &mdash;&nbsp;Нет, ухожу."},
        {"Так бывает. И вот так...И ещё вот так!.. Бывает же???Что поделать.", "Так бывает. И&nbsp;вот так&hellip;И ещё вот так!.. Бывает&nbsp;же???Что поделать."},
        {"Кое от кого, кое на чем, кой у кого, кое с чьим.", "Кое от&nbsp;кого, кое на&nbsp;чем, кой у&nbsp;кого, кое с&nbsp;чьим."},
        {"№ 15Ф, № 34/25", "&#8470;&nbsp;15Ф, &#8470;&nbsp;34/25"},
        {"text...  ... text ...Text", "text&hellip; &hellip; text &hellip;Text"},
        {"... text ...Text", "&hellip; text &hellip;Text"},
        {"то же, сказал бы, думал ли я ли", "то&nbsp;же, сказал&nbsp;бы, думал&nbsp;ли я&nbsp;ли"},
        {"text(text )text", "text (text) text"},
        {"text(text ).", "text (text)."},
        {"text(text )((text)text)", "text (text) ((text) text)"},
        {"броненосец \"Потёмкин\" выдал следующее", "броненосец &laquo;Потёмкин&raquo; выдал следующее"},
        {"броненосец \" Потёмкин", "броненосец &quot; Потёмкин"},
        {"броненосец \"\"Потёмкин\"", "броненосец &laquo;&laquo;Потёмкин&raquo;"},
        {"\"Потёмкин\" броненосец ", "&laquo;Потёмкин&raquo; броненосец"},
        {"I & Sveta", "I&nbsp;&amp; Sveta"},
        {"Dima & Sveta", "Dima &amp; Sveta"},
        {"§32, §IV", "&sect;&nbsp;32, &sect;&nbsp;IV"},
        {"text §, §IV", "text &sect;, &sect;&nbsp;IV"},
        {"&sect; 32", "&sect;&nbsp;32"},
        {"В Пятом §.", "В&nbsp;Пятом &sect;."},
        {"Это и есть секрет - жизнь после смерти.", "Это и&nbsp;есть секрет&nbsp;&mdash; жизнь после смерти."},
        {"5-5", "5-5"},
        {"<p>Это и есть секрет - жизнь после смерти.</p>", "<p>Это и&nbsp;есть секрет&nbsp;&mdash; жизнь после смерти.</p>"},
        {"<p>Когда-нибудь это должно было случиться.</p><p>\"Джедаи\" вернулись.</p>", "<p>Когда-нибудь это должно было случиться.</p><p>&laquo;Джедаи&raquo; вернулись.</p>"},
        {"<b>Книга</b> - это", "<b>Книга</b>&nbsp;&mdash; это"},
        {"", ""},
        {"", ""},
        {"", ""},
        {"", ""},
        {"", ""},
        {"", ""},
        {"", ""},
        {"", ""},


    }
)

func TestTypo(t *testing.T) {
    for _, test := range testerTypo {
        res := Typo(test.In)
        if res != test.Out {
            t.Errorf("%q == %q, NEED: %q", test.In, res, test.Out)
        }
    }
}