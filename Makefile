APP = gmenu

# Output binaries
OUT_MAC = $(APP)-darwin
OUT_WINDOWS = $(APP)-windows

CFLAGS = -O2 -march=native -Wall -Wextra -flto

macos:
	gcc $(CFLAGS) -framework Cocoa main.c search.c platform_darwin.m -o $(OUT_MAC)

windows:
	gcc $(CFLAGS) main.c search.c platform_windows.c -lgdi32 -o $(OUT_WINDOWS)

