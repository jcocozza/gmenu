APP = gmenu

# Windows cross-compiler
CC_WIN = x86_64-w64-mingw32-gcc

# Output binaries
OUT_LINUX = $(APP)-linux-amd64
OUT_WINDOWS = $(APP)-windows-amd64.exe
OUT_MAC_X86 = $(APP)-darwin-amd64
OUT_MAC_ARM = $(APP)-darwin-arm64

# Raylib includes and sources
RAYLIB_SRC = raylib/src
RAYLIB_INC = -I$(RAYLIB_SRC) -I$(RAYLIB_SRC)/external -I$(RAYLIB_SRC)/external/glfw/include
RAYLIB_OBJS = rcore.o rshapes.o rtext.o rmodels.o rtextures.o raudio.o utils.o

.PHONY: all clean linux windows macos-x86_64 macos-arm64

all:
	@echo "Please specify a target: linux, windows, macos-x86_64, macos-arm64"

clean:
	rm -f $(OUT_LINUX) $(OUT_WINDOWS) $(OUT_MAC_X86) $(OUT_MAC_ARM)
	rm -f libraylib.a $(RAYLIB_SRC)/*.o

# Build raylib as static library
raylib-static:
	cd $(RAYLIB_SRC) && \
	gcc -DPLATFORM_DESKTOP -c rcore.c rshapes.c rtext.c rmodels.c rtextures.c raudio.c utils.c \
	    -I. -Iexternal -Iexternal/glfw/include
	ar rcs libraylib.a $(RAYLIB_SRC)/*.o

linux: raylib-static
	gcc -O2 -flto -o $(OUT_LINUX) main.c \
	    $(RAYLIB_INC) libraylib.a -lGL -lm -ldl -lpthread -lrt -lX11

windows: raylib-static
	$(CC_WIN) -O2 -o $(OUT_WINDOWS) main.c \
	    $(RAYLIB_INC) libraylib.a -lopengl32 -lgdi32 -lwinmm -lshell32

macos-x86_64:
	gcc -O2 -target x86_64-apple-macos12 -ObjC main.c window_darwin.m -o $(OUT_MAC_X86) \
	    -I/opt/homebrew/opt/raylib/include \
	    -L/opt/homebrew/opt/raylib/lib -lraylib \
	    -framework Cocoa -framework IOKit -framework CoreVideo -framework OpenGL

macos-arm64:
	gcc -O2 -target arm64-apple-macos12 -ObjC main.c window_darwin.m -o $(OUT_MAC_ARM) \
	    -I/opt/homebrew/opt/raylib/include \
	    -L/opt/homebrew/opt/raylib/lib -lraylib \
	    -framework Cocoa -framework IOKit -framework CoreVideo -framework OpenGL

