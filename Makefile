APP = gmenu

# Output binaries
OUT_MAC = $(APP)-darwin
OUT_WINDOWS = $(APP)-windows

windows:
	gcc -municode main.c search.c platform_windows.c -lgdi32 -o $(OUT_WINDOWS)

macos:
	gcc -framework Cocoa main.c search.c platform_darwin.m -o $(OUT_MAC)
