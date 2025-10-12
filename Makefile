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

# I do not understand this at all. No clue how I got this working.
# It would probably be much easier if I wasn't trying to do with on an ARM Windows VM...
#
# this assumes you have libraylib.a file and a clone of raylib in the root of gmenu.
#
# generate the libraylib.a file with:
# 1. gcc -c src/rcore.c src/rglfw.c src/rshapes.c src/rtext.c src/rmodels.c src/rtextures.c src/raudio.c src/utils.c -I"src" -I"src/external" -I"src/external/glfw/include" -DPLATFORM_DESKTOP
# 2. ar rcs libraylib.a rcore.o rglfw.o rshapes.o rtext.o rmodels.o rtextures.o raudio.o utils.o
#
# Then move the libraylib.a file to the root of gmenu.
#
# then use `mingw32-make windows` or just run the following gcc command:
windows:
	gcc main.c -I"raylib\src" -I"raylib\src\external" -I"raylib\src\external\glfw\include" -L. -lraylib -lopengl32 -lgdi32 -lwinmm -lshell32 -o gmenu.exe


linux:
	gcc -flto -O3 -DNDEBUG -march=native main.c -lraylib -lGL -lm -lpthread -ldl -o $(OUT)
	strip $(OUT)


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
	strip $(OUT)
