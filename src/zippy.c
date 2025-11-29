#include <stdio.h>

#include "lex.h"

int main() {
    printf("Hello, %s!\n", token_type_to_str(T_BANG));
}
