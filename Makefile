.PHONY: all clean

PLATFORM ?= $(shell uname -s)
BUILD_DIR ?= build
THIRDPARTY_DIR ?= thirdparty
SRC_DIR ?= src

THIRDPARTY_INCLUDE = $(THIRDPARTY_DIR)/include
THIRDPARTY_LIBS = $(THIRDPARTY_DIR)/lib
CC=clang
CFLAGS = -O3 -Wall -std=c23 -I$(THIRDPARTY_INCLUDE)
AR = ar -crs
LDFLAGS = -L$(THIRDPARTY_LIBS) -lutf8proc

ZIPPY_EXE = $(BUILD_DIR)/zippy
EXE_OBJ = $(BUILD_DIR)/zippy.o
ZIPPY_STATIC_LIB = $(BUILD_DIR)/libzippy.a

SRCS = $(filter-out $(SRC_DIR)/zippy.c, $(wildcard $(SRC_DIR)/*.c))
OBJS = $(sort $(SRCS:$(SRC_DIR)/%.c=$(BUILD_DIR)/%.o))

all: $(ZIPPY_EXE) $(ZIPPY_STATIC_LIB)

$(BUILD_DIR)/%.o: $(SRC_DIR)/%.c
	@mkdir -p -- $(dir $@)
	$(CC) $(CFLAGS) -o $@ -c $^

$(ZIPPY_STATIC_LIB): $(OBJS)
	$(AR) $(ZIPPY_STATIC_LIB) $^

$(ZIPPY_EXE): $(EXE_OBJ) $(ZIPPY_STATIC_LIB)
	@mkdir -p -- $(dir $@)
	$(CC) $(CFLAGS) -o $(ZIPPY_EXE) $^ $(LDFLAGS)

clean:
	@echo "Removing build artifacts..."
	rm -rdf $(BUILD_DIR)
	@echo "Removed!"
