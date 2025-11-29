// TODO: add more zippy status codes
#include <stdio.h>
#include <stdlib.h>

#include "common.h"
#include "lex.h"

const char *token_type_to_str(token_type type) {
    switch (type) {
    case T_INVALID:
        return "INVALID";
    case T_EOF:
        return "EOF";
    case T_BREAK:
        return "BREAK";
    case T_CASE:
        return "CASE";
    case T_CONST:
        return "CONST";
    case T_CONTINUE:
        return "CONTINUE";
    case T_DEFAULT:
        return "DEFAULT";
    case T_ELSE:
        return "ELSE";
    case T_ENUM:
        return "ENUM";
    case T_FN:
        return "FN";
    case T_FOR:
        return "FOR";
    case T_IMPORT:
        return "IMPORT";
    case T_INTERFACE:
        return "INTERFACE";
    case T_IF:
        return "IF";
    case T_IN:
        return "IN";
    case T_MAP:
        return "MAP";
    case T_RETURN:
        return "RETURN";
    case T_STR:
        return "STR";
    case T_STRUCT:
        return "STRUCT";
    case T_SWITCH:
        return "SWITCH";
    case T_TYPE:
        return "TYPE";
    case T_VAR:
        return "VAR";
    case T_WEAK:
        return "WEAK";
    case T_PLUS:
        return "PLUS";
    case T_MINUS:
        return "MINUS";
    case T_STAR:
        return "STAR";
    case T_SLASH:
        return "SLASH";
    case T_PERCENT:
        return "PERCENT";
    case T_AMPERSAND:
        return "AMPERSAND";
    case T_PIPE:
        return "PIPE";
    case T_TILDE:
        return "TILDE";
    case T_CARET:
        return "CARET";
    case T_LEFT_SHIFT:
        return "LEFT_SHIFT";
    case T_RIGHT_SHIFT:
        return "RIGHT_SHIFT";
    case T_PLUS_EQ:
        return "PLUS_EQ";
    case T_MINUS_EQ:
        return "MINUS_EQ";
    case T_STAR_EQ:
        return "STAR_EQ";
    case T_SLASH_EQ:
        return "SLASH_EQ";
    case T_PERCENT_EQ:
        return "PERCENT_EQ";
    case T_AND_EQ:
        return "AND_EQ";
    case T_OR_EQ:
        return "OR_EQ";
    case T_NOT_EQ_BIT:
        return "TILDE_EQ";
    case T_LSHIFT_EQ:
        return "LSHIFT_EQ";
    case T_RSHIFT_EQ:
        return "RSHIFT_EQ";
    case T_WALRUS:
        return "WALRUS_EQ";
    case T_EQUAL:
        return "EQUAL";
    case T_AND:
        return "AND_LOGICAL";
    case T_OR:
        return "OR_LOGICAL";
    case T_BANG:
        return "BANG";
    case T_BANG_EQ:
        return "NOT_EQUAL";
    case T_EQUAL_EQUAL:
        return "EQUAL_EQUAL";
    case T_LESS:
        return "LESS";
    case T_GREATER:
        return "GREATER";
    case T_LESS_EQ:
        return "LESS_EQ";
    case T_GREATER_EQ:
        return "GREATER_EQ";
    case T_PLUS_PLUS:
        return "INCREMENT";
    case T_MINUS_MINUS:
        return "DECREMENT";
    case T_QUESTION:
        return "QUESTION";
    case T_LEFT_PAREN:
        return "LEFT_PAREN";
    case T_RIGHT_PAREN:
        return "RIGHT_PAREN";
    case T_LEFT_BRACKET:
        return "LEFT_BRACKET";
    case T_RIGHT_BRACKET:
        return "RIGHT_BRACKET";
    case T_LEFT_BRACE:
        return "LEFT_BRACE";
    case T_RIGHT_BRACE:
        return "RIGHT_BRACE";
    case T_COMMA:
        return "COMMA";
    case T_SEMICOLON:
        return "SEMICOLON";
    case T_COLON:
        return "COLON";
    case T_DOUBLE_COLON:
        return "DOUBLE_COLON";
    case T_DOT:
        return "DOT";
    case T_DOT_DOT:
        return "DOT_DOT";
    case T_IDENTIFIER:
        return "IDENTIFIER";
    case T_INTEGER:
        return "INTEGER";
    case T_REAL:
        return "REAL_NUMBER";
    case T_CHAR:
        return "CHARACTER";
    case T_STRING:
        return "STRING";
    default:
        return "UNKNOWN";
    }
}

zippy_status lexer_destroy(zippy_lexer *zl) {
    if (zl->text != nullptr) {
        free(zl->text);
    } else {
        return err;
    }
    return ok;
}

zippy_status lexer_init(zippy_lexer *zl, const char *module_name) {
    if (zl->is_module) {
        FILE *module = fopen(module_name, "rb");
        if (module == nullptr) {
            return err;
        }

        fseek(module, 0, SEEK_END);
        long txt_len = ftell(module);
        rewind(module);

        zl->text = (u8 *)malloc(txt_len + 1);
        if ((long)fread(zl->text, 1, txt_len, module) != txt_len) {
            return err;
        }

        zl->text_len = txt_len;
        zl->text[txt_len] = '\0';
        fclose(module);
    } else {
        // TODO: read from a source string (script)
    }

    zl->cursor = 0;
    zl->line_no = 1;
    zl->column_no = 1;
    zl->token = INVALID_TOKEN;
    zl->last_token = zl->token;
    zl->is_eof = false;
    return ok;
}

zippy_internal void next_rune(zippy_lexer *zl) {
    if (zl->cursor < zl->text_len) {
        rune c = (rune)zl->text[zl->cursor];
        if (c == eof) {
            // TODO: Place and error here
            exit(-1);
        } else if ((c & 0x80) != 0) { // non-ASCII character
            usize width = utf8_decode(zl->text + zl->cursor,
                                      zl->text_len - zl->cursor, &c);
            zl->cursor += width;
            if (c == invalid_rune && width == 1) {
                // TODO: Place error here
            } else if (c == bom_rune && zl->cursor > 0) {
                // TODO: Place error here
            }
        } else {
            zl->cursor++;
        }
        zl->curr_rune = c;
    } else {
        zl->is_eof = true;
    }
}
