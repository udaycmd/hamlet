#include "common.h"

usize utf8_decode(const u8 *str, usize len, rune *cp_out) {
    usize width = 1;
    rune codepoint = invalid_rune;

    if (len > 1) {
        u8 b0 = str[0];

        if ((b0 & 0xE0) == 0xC0) {
            codepoint = ((b0 & 0x1F) << 6) | ((str[1] & 0x3F));
            width = 2;
        }

        if ((b0 & 0xF0) == 0xE0) {
            if (len < 3)
                goto end;

            codepoint = ((b0 & 0x0F) << 12) | ((str[1] & 0x3F) << 6) |
                        ((str[2] & 0x3F));
            width = 3;
        }

        if ((b0 & 0xF8) == 0xF0) {
            if (len < 4)
                goto end;

            codepoint = ((b0 & 0x07) << 18) | ((str[1] & 0x3F) << 12) |
                        ((str[2] & 0x3F) << 6) | ((str[3] & 0x3F));
            width = 4;
        }
    } else {
        width = len;
    }

    if (codepoint > max_rune || (0xD800u <= codepoint && codepoint <= 0xDFFFu))
        codepoint = invalid_rune;

end:
    if (cp_out != nullptr)
        *cp_out = codepoint;
    return width;
}

bool is_letter(rune r) {
    if (r < 0x80) {
        if (r == '_') {
            return true;
        }
        return ((r | 0x20) - 0x61) < 26;
    }

    switch (utf8proc_category(r)) {
    case UTF8PROC_CATEGORY_LU:
    case UTF8PROC_CATEGORY_LL:
    case UTF8PROC_CATEGORY_LT:
    case UTF8PROC_CATEGORY_LM:
    case UTF8PROC_CATEGORY_LO:
        return true;
    default:
        return false;
    }
}

bool is_digit(rune r) {
    if (r < 0x80) {
        return (r - '0') < 10;
    }
    return utf8proc_category(r) == UTF8PROC_CATEGORY_ND;
}

bool is_whitespace(rune r) {
    switch (r) {
    case ' ':
    case '\t':
    case '\n':
    case '\r':
        return true;
    }
    return false;
}
