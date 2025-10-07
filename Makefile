build:
	gcc -g main.c -lraylib -lGL -lm -lpthread -ldl -o gmenu

build-macos:
	gcc -g main.c     -I/opt/homebrew/opt/raylib/include     -L/opt/homebrew/opt/raylib/lib     -lraylib     -framework Cocoa     -framework IOKit     -framework CoreVideo     -framework OpenGL     -o gmenu
