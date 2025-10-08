build:
	gcc -g main.c -lraylib -lGL -lm -lpthread -ldl -o gmenu

macos:
	gcc -o gmenu -ObjC main.c window_darwin.m \
		-I/opt/homebrew/opt/raylib/include -L/opt/homebrew/opt/raylib/lib \
		-lraylib \
		-framework Cocoa \
		-framework IOKit \
		-framework CoreVideo \
		-framework OpenGL
