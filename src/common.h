#include <stdint.h>
#include <uchar.h>

#include "utf8proc.h"

#define zippy_internal static

// unicode(utf-8) definations

typedef char32_t rune;
typedef uint8_t u8;
typedef uint64_t u64;
typedef size_t usize;

#define eof (uint8_t)0
#define invalid_rune (rune)0xfffd
#define bom_rune (rune)0xfeff
#define max_rune (rune)0x10FFFF

usize utf8_decode(const u8 *str, usize len, rune *cp_out);
bool is_letter(rune r);
bool is_digit(rune r);
bool is_whitespace(rune r);
