# linux:
#		gcc -O3 -DNDEBUG main.c -lraylib -lGL -lm -lpthread -ldl -o gmenu

# Paths to Raylib installed via Homebrew
RAYLIB_INC = /opt/homebrew/opt/raylib/include
RAYLIB_LIB = /opt/homebrew/opt/raylib/lib

# Common framework flags for macOS
FRAMEWORKS = -framework Cocoa -framework IOKit -framework CoreVideo -framework OpenGL

# Source files
SRC_C = main.c
SRC_OBJC = window_darwin.m

OUT = gmenu

all: debug

clean:
	rm -f $(OUT)

# these are all for macos
debug:
	gcc -g -O0 -Wall -ObjC $(SRC_C) $(SRC_OBJC) \
		-I$(RAYLIB_INC) -L$(RAYLIB_LIB) \
		-lraylib $(FRAMEWORKS) \
		-o $(OUT)

release:
	gcc -O3 -DNDEBUG -ObjC $(SRC_C) $(SRC_OBJC) \
		-I$(RAYLIB_INC) -L$(RAYLIB_LIB) \
		-lraylib $(FRAMEWORKS) \
		-o $(OUT)

# I don't really understand this, but supposedly it should improve startup times
release_static:
	gcc -flto -O3 -DNDEBUG -march=native -ObjC $(SRC_C) $(SRC_OBJC) \
		-I$(RAYLIB_INC) $(RAYLIB_LIB)/libraylib.a \
		$(FRAMEWORKS) \
		-o $(OUT)
	strip $(OUT)_static
