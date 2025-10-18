APP = gmenu

CFLAGS = -O2 -march=native -Wall -Wextra -flto

macos:
	gcc $(CFLAGS) -framework Cocoa main.c search.c platform_darwin.m -o $(APP)

windows:
	gcc $(CFLAGS) main.c search.c platform_windows.c -lgdi32 -o $(APP)

