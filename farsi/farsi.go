
package farsi

import (
    _ "regexp"
    "unicode"
)

// Farsi struct to mimic PHP class
type Farsi struct {
    pChars   map[rune][3]string
    tahoma   map[rune][3]string
    normal   map[rune][3]string
    mpChars  map[rune]bool
    ignorelist map[rune]bool
    openClose map[rune]bool
    enChars  map[rune]bool
}

// Constructor
func NewFarsi() *Farsi {
    return &Farsi{
        pChars: map[rune][3]string{
            'آ': {"ﺂ", "ﺂ", "آ"},
            'ا': {"ﺎ", "ﺎ", "ا"},
            'ب': {"ﺐ", "ﺒ", "ﺑ"},
            'پ': {"ﭗ", "ﭙ", "ﭘ"},
            'ت': {"ﺖ", "ﺘ", "ﺗ"},
            'ث': {"ﺚ", "ﺜ", "ﺛ"},
            'ج': {"ﺞ", "ﺠ", "ﺟ"},
            'چ': {"ﭻ", "ﭽ", "ﭼ"},
            'ح': {"ﺢ", "ﺤ", "ﺣ"},
            'خ': {"ﺦ", "ﺨ", "ﺧ"},
            'د': {"ﺪ", "ﺪ", "ﺩ"},
            'ذ': {"ﺬ", "ﺬ", "ﺫ"},
            'ر': {"ﺮ", "ﺮ", "ﺭ"},
            'ز': {"ﺰ", "ﺰ", "ﺯ"},
            'ژ': {"ﮋ", "ﮋ", "ﮊ"},
            'س': {"ﺲ", "ﺴ", "ﺳ"},
            'ش': {"ﺶ", "ﺸ", "ﺷ"},
            'ص': {"ﺺ", "ﺼ", "ﺻ"},
            'ض': {"ﺾ", "ﻀ", "ﺿ"},
            'ط': {"ﻂ", "ﻄ", "ﻃ"},
            'ظ': {"ﻆ", "ﻈ", "ﻇ"},
            'ع': {"ﻊ", "ﻌ", "ﻋ"},
            'غ': {"ﻎ", "ﻐ", "ﻏ"},
            'ف': {"ﻒ", "ﻔ", "ﻓ"},
            'ق': {"ﻖ", "ﻘ", "ﻗ"},
            'ک': {"ﻚ", "ﻜ", "ﻛ"},
            'گ': {"ﮓ", "ﮕ", "ﮔ"},
            'ل': {"ﻞ", "ﻠ", "ﻟ"},
            'م': {"ﻢ", "ﻤ", "ﻣ"},
            'ن': {"ﻦ", "ﻨ", "ﻧ"},
            'و': {"ﻮ", "ﻮ", "ﻭ"},
            'ی': {"ﯽ", "ﯿ", "ﯾ"},
            'ك': {"ﻚ", "ﻜ", "ﻛ"},
            'ي': {"ﻲ", "ﻴ", "ﻳ"},
            'أ': {"ﺄ", "ﺄ", "ﺃ"},
            'ؤ': {"ﺆ", "ﺆ", "ﺅ"},
            'إ': {"ﺈ", "ﺈ", "ﺇ"},
            'ئ': {"ﺊ", "ﺌ", "ﺋ"},
            'ة': {"ﺔ", "ﺘ", "ﺗ"},
        },
        tahoma: map[rune][3]string{
            'ه': {"ﮫ", "ﮭ", "ﮬ"},
        },
        normal: map[rune][3]string{
            'ه': {"ﻪ", "ﻬ", "ﻫ"},
        },
        mpChars: map[rune]bool{
            'آ': true, 'ا': true, 'د': true, 'ذ': true, 'ر': true,
            'ز': true, 'ژ': true, 'و': true, 'أ': true, 'إ': true, 'ؤ': true,
        },
        ignorelist: map[rune]bool{
            'ٌ': true, 'ٍ': true, 'ً': true, 'ُ': true, 'ِ': true,
            'َ': true, 'ّ': true, 'ٓ': true, 'ٰ': true, 'ٔ': true, 'ﹶ': true,
            'ﹺ': true, 'ﹸ': true, 'ﹼ': true, 'ﹾ': true, 'ﹴ': true, 'ﹰ': true,
            'ﱞ': true, 'ﱟ': true, 'ﱠ': true, 'ﱡ': true, 'ﱢ': true, 'ﱣ': true,
        },
        openClose: map[rune]bool{
            '>': true, ')': true, '}': true, ']': true,
            '<': true, '(': true, '{': true, '[': true,
        },
        enChars: map[rune]bool{
            'a': true, 'b': true, 'c': true, 'd': true, 'e': true,
            'f': true, 'g': true, 'h': true, 'i': true, 'j': true,
            'k': true, 'l': true, 'm': true, 'n': true, 'o': true,
            'p': true, 'q': true, 'r': true, 's': true, 't': true,
            'u': true, 'v': true, 'w': true, 'x': true, 'y': true, 'z': true,
        },
    }
}

// utf8_strlen: count UTF-8 characters
func (f *Farsi) utf8Len(s string) int {
    return len([]rune(s))
}

// fa_number: convert Latin digits to Persian/Arabic digits
func (f *Farsi) faNumber(num string) string {
    AF := map[rune]rune{
        '0': '٠', '1': '١', '2': '٢', '3': '٣', '4': '۴',
        '5': '۵', '6': '۶', '7': '٧', '8': '٨', '9': '٩',
    }
    var result []rune
    for _, r := range num {
        if mapped, ok := AF[r]; ok {
            result = append(result, mapped)
        } else {
            result = append(result, r)
        }
    }
    return string(result)
}

// PersiaText: main function to process Persian text
func (f *Farsi) PersiaText(str string, z string, method string, farsiNumber bool) string {
    var output, enStr, num string
    var runWay []string // for debugging (optional)

    // Merge pChars with tahoma or normal based on method
    pChars := make(map[rune][3]string)
    for k, v := range f.pChars {
        pChars[k] = v
    }
    if method == "tahoma" {
        for k, v := range f.tahoma {
            pChars[k] = v
        }
    } else {
        for k, v := range f.normal {
            pChars[k] = v
        }
    }

    runes := []rune(str)
    strLen := len(runes)

    for i := 0; i < strLen; i++ {
        var str1, strNext, strBack rune
        str1 = runes[i]

        // Reset runWay per iteration (optional for debug)
        runWay = nil

        // Determine strNext and strBack safely
        if i+1 < strLen && f.ignorelist[runes[i+1]] {
            if i+2 < strLen {
                strNext = runes[i+2]
            } else {
                strNext = 0
            }
            if i == 2 {
                if i >= 2 {
                    strBack = runes[i-2]
                }
            } else {
                if i > 0 {
                    strBack = runes[i-1]
                }
            }
        } else if !(i > 0 && f.ignorelist[runes[i-1]]) {
            if i+1 < strLen {
                strNext = runes[i+1]
            }
            if i > 0 {
                strBack = runes[i-1]
            }
        } else {
            if i+1 < strLen {
                strNext = runes[i+1]
            } else {
                if i > 0 {
                    strNext = runes[i-1]
                }
            }
            if i >= 2 {
                strBack = runes[i-2]
            }
        }

        if !f.ignorelist[str1] {
            if chars, exists := pChars[str1]; exists {
                if strBack == 0 || strBack == ' ' || !pCharsExists(pChars, strBack) {
                    if !pCharsExists(pChars, strBack) && !pCharsExists(pChars, strNext) {
                        output = string(str1) + output
                    } else {
                        output = chars[2] + output
                    }
                    continue
                } else if pCharsExists(pChars, strNext) && pCharsExists(pChars, strBack) {
                    if f.mpChars[strBack] && pCharsExists(pChars, strNext) {
                        output = chars[2] + output
                    } else {
                        output = chars[1] + output
                    }
                    continue
                } else if pCharsExists(pChars, strBack) && !pCharsExists(pChars, strNext) {
                    if f.mpChars[strBack] {
                        output = string(str1) + output
                    } else {
                        output = chars[0] + output
                    }
                    continue
                }
            } else if z == "fa" {
                numbers := "٠١٢٣٤٥٦٧٨٩۴۵۶0123456789"

                // Flip brackets and symbols
                switch str1 {
                case ')':
                    str1 = '('
                case '(':
                    str1 = ')'
                case '}':
                    str1 = '{'
                case '{':
                    str1 = '}'
                case ']':
                    str1 = '['
                case '[':
                    str1 = ']'
                case '>':
                    str1 = '<'
                case '<':
                    str1 = '>'
                }

                // Handle numbers
                // BUG in numbers persian
                // if isDigit(str1, numbers) {
                //     if farsiNumber {
                //         num += f.faNumber(string(str1))
                //         runWay = append(runWay, "1")
                //     } else {
                //         num += string(str1)
                //         runWay = append(runWay, "2")
                //     }
                //     str1 = 0
                // }

                // Check if next is not number
                if !isDigit(strNext, numbers) {
                    if (isEnglishChar(str1) || ((str1 == ' ' || str1 == '.') && enStr != "" && !pCharsExists(pChars, strNext))) {
                        enStr += string(str1) + num
                        str1 = 0
                        runWay = append(runWay, "3")
                    } else {
                        if enStr != "" {
                            if i+1 == strLen {
                                runWay = append(runWay, "3.5")
                                str1 = 0
                                num = string(str1) + num // reassign?
                            } else {
                                enStr += string(str1) + num
                                runWay = append(runWay, "4")
                            }
                        } else {
                            num = string(str1) + num
                            runWay = append(runWay, "5")
                        }
                    }
                    num = ""
                }

                // Handle English text and context
                if enStr != "" || (str1 != 0 && i == 0 && (!pCharsExists(pChars, strNext) && strNext != ' ')) {
                    if !pCharsExists(pChars, str1) {
                        if !pCharsExists(pChars, strNext) && strNext != ' ' && !f.openClose[strNext] {
                            enStr += string(str1)
                            runWay = append(runWay, "6")
                        } else {
                            if i+2 < strLen && f.enChars[unicode.ToLower(runes[i+2])] {
                                enStr += string(str1)
                                runWay = append(runWay, "7")
                            } else if strNext == ' ' && i+2 < strLen && (isDigit(runes[i+2], numbers) || f.enChars[unicode.ToLower(runes[i+2])]) {
                                enStr += string(str1)
                                runWay = append(runWay, "8")
                            } else {
                                output = enStr + output
                                enStr = ""
                                runWay = append(runWay, "9")
                            }
                        }
                    } else {
                        if num != "" {
                            enStr += num
                            runWay = append(runWay, "10")
                        } else {
                            output = enStr + string(str1) + output
                            enStr = ""
                            runWay = append(runWay, "11")
                        }
                    }
                } else {
                    if isDigit(str1, numbers) && strNext == '.' && i+2 < strLen && isDigit(runes[i+2], numbers) {
                        enStr = string(str1)
                        runWay = append(runWay, "12")
                    } else {
                        output = string(str1) + output
                        runWay = append(runWay, "14")
                    }
                }
            } else {
                // Handle punctuation and spacing
                if str1 == '،' || str1 == '؟' || str1 == 'ء' ||
                    (pCharsExists(pChars, strNext) && pCharsExists(pChars, strBack)) ||
                    (str1 == ' ' && pCharsExists(pChars, strBack)) ||
                    (str1 == ' ' && pCharsExists(pChars, strNext)) {
                    output = string(str1) + output
                } else {
                    enStr += string(str1)
                    if pCharsExists(pChars, strNext) || strNext == 0 {
                        output = enStr + output
                        enStr = ""
                    }
                }
            }
        } else {
            output = string(str1) + output
        }
    }

    if enStr != "" {
        output = enStr + output
    }

    return output
}

// Helper functions
func pCharsExists(p map[rune][3]string, r rune) bool {
    _, exists := p[r]
    return exists
}

func isDigit(r rune, list string) bool {
    for _, ch := range list {
        if ch == r {
            return true
        }
    }
    return false
}

func isEnglishChar(r rune) bool {
    return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == ' ' || r == '.'
}
