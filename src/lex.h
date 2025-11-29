#pragma once

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#include "common.h"
#include "zippy.h"

typedef enum ttype {
    T_INVALID = 0,
    T_EOF,

    // - Keywords -
    T_BREAK,
    T_CASE,
    T_CONST,
    T_CONTINUE,
    T_DEFAULT,
    T_ELSE,
    T_ENUM,
    T_FN,
    T_FOR,
    T_IMPORT,
    T_INTERFACE,
    T_IF,
    T_IN,
    T_MAP,
    T_RETURN,
    T_STR,
    T_STRUCT,
    T_SWITCH,
    T_TYPE,
    T_VAR,
    T_WEAK,

    // - Operators -
    T_PLUS,
    T_MINUS,
    T_STAR,
    T_SLASH,
    T_PERCENT,
    T_AMPERSAND,
    T_PIPE,
    T_TILDE,
    T_CARET,
    T_LEFT_SHIFT,
    T_RIGHT_SHIFT,
    T_PLUS_EQ,
    T_MINUS_EQ,
    T_STAR_EQ,
    T_SLASH_EQ,
    T_PERCENT_EQ,
    T_AND_EQ,
    T_OR_EQ,
    T_NOT_EQ_BIT,
    T_LSHIFT_EQ,
    T_RSHIFT_EQ,
    T_WALRUS,
    T_EQUAL,
    T_AND,
    T_OR,
    T_BANG,
    T_BANG_EQ,
    T_EQUAL_EQUAL,
    T_LESS,
    T_GREATER,
    T_LESS_EQ,
    T_GREATER_EQ,
    T_PLUS_PLUS,
    T_MINUS_MINUS,
    T_QUESTION,

    // - Punctuations -
    T_LEFT_PAREN,
    T_RIGHT_PAREN,
    T_LEFT_BRACKET,
    T_RIGHT_BRACKET,
    T_LEFT_BRACE,
    T_RIGHT_BRACE,
    T_COMMA,
    T_SEMICOLON,
    T_COLON,
    T_DOUBLE_COLON,
    T_DOT,
    T_DOT_DOT,

    // - Literals -
    T_IDENTIFIER,
    T_INTEGER,
    T_REAL,
    T_CHAR,
    T_STRING,
} token_type;

typedef struct token {
    u8 *start;
    usize length;
    token_type type;
    u64 line;
    u64 col;
} zippy_token;

typedef struct lexer {
    u8 *text;
    usize text_len;
    bool is_module;
    const char *module_name;
    u64 line_no;
    u64 column_no;
    u64 cursor;
    rune curr_rune;
    bool is_eof;
    zippy_token token, last_token;
} zippy_lexer;

const char *token_type_to_str(token_type type);
zippy_status lexer_init(zippy_lexer *zl, const char *module_name);
zippy_status lexer_destroy(zippy_lexer *zl);
void next(zippy_lexer *zl);

static const zippy_token INVALID_TOKEN = {
    .col = 0, .length = 0, .line = 0, .start = nullptr, .type = T_INVALID};
