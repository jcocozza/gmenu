#include <windows.h>
#include <stdio.h>
#include "platform.h"

static int height = 30;
static HWND hwnd = NULL;
static HINSTANCE h_instance_global = NULL;
static gmenu_keypress_t last_kp = {0}; // this will be updated to reflect the last keypress
int close = 0;

LRESULT CALLBACK WndProc(HWND hwnd, UINT msg, WPARAM w_param, LPARAM l_param);

LRESULT CALLBACK WndProc(HWND hwnd, UINT msg, WPARAM w_param, LPARAM l_param) {
    switch (msg) {
        case WM_DESTROY:
            PostQuitMessage(0);
            return 0;
	case WM_CHAR: {
		char ch = (char)w_param;
    		if (ch >= 32 && ch <= 126) { // printable ASCII
        		last_kp.k = KEY_CHAR;
        		last_kp.c = ch;
    		}
    		break;
        }
        case WM_KEYDOWN: {
            switch (w_param) {
                case VK_LEFT:   last_kp.k = KEY_LEFT;   break;
                case VK_RIGHT:  last_kp.k = KEY_RIGHT;  break;
                case VK_RETURN: last_kp.k = KEY_ENTER;  break;
                case VK_BACK: last_kp.k = KEY_BACKSPACE;  break;
                case VK_ESCAPE: last_kp.k = KEY_ESC; close = 1; break;
                default: break;
            }
            break;
        }
    }
    return DefWindowProc(hwnd, msg, w_param, l_param);
}

void init() {
    h_instance_global = GetModuleHandle(NULL);

    const char CLASS_NAME[] = "GMenuOverlay";
    WNDCLASS wc = {0};

    wc.lpfnWndProc = WndProc;
    wc.hInstance = h_instance_global;
    wc.lpszClassName = CLASS_NAME;
    wc.hbrBackground = (HBRUSH)GetStockObject(NULL_BRUSH);
    wc.hCursor = LoadCursor(NULL, IDC_ARROW);
    wc.hbrBackground = CreateSolidBrush(RGB(255,255,255));

    if (!RegisterClass(&wc)) return;

    int screen_width  = GetSystemMetrics(SM_CXSCREEN);

    hwnd = CreateWindowEx(
        WS_EX_TOPMOST | WS_EX_LAYERED | WS_EX_TOOLWINDOW,
        CLASS_NAME, "Overlay",
        WS_POPUP, 0, 0, screen_width, height,
        NULL, NULL, h_instance_global, NULL
    );

    if (!hwnd) return;

    SetLayeredWindowAttributes(hwnd, 0, 255, LWA_ALPHA);
    ShowWindow(hwnd, SW_SHOW);
}

void teardown() {
    if (hwnd) {
        DestroyWindow(hwnd);
        hwnd = NULL;
    }
}

int should_close() {
    return close;
}

gmenu_keypress_t get_key_press() {
    MSG msg;
    while (PeekMessage(&msg, NULL, 0, 0, PM_REMOVE)) {
        TranslateMessage(&msg);
        DispatchMessage(&msg);
    }

    gmenu_keypress_t kp = last_kp;
    last_kp.k = KEY_NONE;
    last_kp.c = 0;
    return kp;
}

int text_width(char *txt) {
    	HDC hdc = GetDC(hwnd);
	SIZE size;
	GetTextExtentPoint32A(hdc, txt, strlen(txt), &size);
	return size.cx;
}

int screen_width() {
    return GetSystemMetrics(SM_CXSCREEN);
}

// populated and released by begin and end draw
HDC hdc;
void begin_draw() {
    hdc = GetDC(hwnd);
}

void end_draw() {
    ReleaseDC(hwnd, hdc);
}

// TODO: figure out how to do the coloring properly
void draw_text(char *txt, int x, int y, gmenu_color_t c) {
	if (hdc == NULL) { return; }
	COLORREF old_color; 
	int old_bk_mode = SetBkMode(hdc, TRANSPARENT);
	switch (c) {
		case GMENU_RED:	 
			old_color = SetTextColor(hdc, RGB(255, 0, 0)); 
			break;
		case GMENU_BLACK: 
			old_color = SetTextColor(hdc, RGB(0, 0, 0)); 
			break;
	}
	TextOutA(hdc, x, y, txt, strlen(txt));
}

void clear_screen() {
    if (hdc == NULL) { return; }
    RECT rect;
    GetClientRect(hwnd, &rect);
    HBRUSH brush = CreateSolidBrush(RGB(255,255,255));
    FillRect(hdc, &rect, brush);
    DeleteObject(brush);
}
