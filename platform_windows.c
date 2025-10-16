#include <windows.h>
#include <stdio.h>
#include "platform.h"

static HWND hwnd = NULL;
static HINSTANCE h_instance_global = NULL;

LRESULT CALLBACK WndProc(HWND hwnd, UINT msg, WPARAM w_param, LPARAM l_param);

LRESULT CALLBACK WndProc(HWND hwnd, UINT msg, WPARAM w_param, LPARAM l_param) {
    switch (msg) {
        case WM_DESTROY:
            PostQuitMessage(0);
            return 0;
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

    hwnd = CreateWindowEx(
        WS_EX_TOPMOST | WS_EX_LAYERED | WS_EX_TOOLWINDOW,
        CLASS_NAME, "Overlay",
        WS_POPUP, 100, 100, 400, 200,
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
    return 0;
}

#define KEY_COUNT 256

int is_key_pressed_once(int vk) {
    static bool prev_state[KEY_COUNT] = {0};
    SHORT state = GetAsyncKeyState(vk);
    bool pressed_now = (state & 0x8000) != 0;

    bool pressed_once = pressed_now && !prev_state[vk];
    prev_state[vk] = pressed_now;
    return pressed_once;
}

gmenu_keypress_t get_key_press() {
    gmenu_keypress_t kp = {};

    if (is_key_pressed_once(VK_LEFT))
        kp.k = KEY_LEFT;
    else if (is_key_pressed_once(VK_RIGHT))
        kp.k = KEY_RIGHT;
    else if (is_key_pressed_once(VK_RETURN))
        kp.k = KEY_ENTER;
    else if (is_key_pressed_once(VK_ESCAPE))
        kp.k = KEY_ESC;
    else {
        for (int c = 'A'; c <= 'Z'; c++) {
            if (is_key_pressed_once(c)) {
                kp.k = KEY_CHAR;
                kp.c = (char)c;
                break;
            }
        }
        for (int c = '0'; c <= '9'; c++) {
            if (is_key_pressed_once(c)) {
                kp.k = KEY_CHAR;
                kp.c = (char)c;
                break;
            }
        }
    }

    return kp;
}

void draw(char *user_prompt, char *user_input, search_results_t *results, int result_offset, int selected_result) {
    if (!hwnd) {
        return;
    }
    HDC hdc = GetDC(hwnd);
    RECT rect = {50, 50, 300, 150};
    // HBRUSH brush = CreateSolidBrush(RGB(255,255,200));
    HPEN pen = CreatePen(PS_SOLID, 2, RGB(0,0,0));
    // HGDIOBJ old_brush = SelectObject(hdc, brush);
    HGDIOBJ old_pen = SelectObject(hdc, pen);

    //Rectangle(hdc, rect.left, rect.top, rect.right, rect.bottom);

    SetBkMode(hdc, TRANSPARENT);
    SetTextColor(hdc, RGB(0,0,0));
    DrawTextA(hdc, user_input, -1, &rect, DT_CENTER | DT_VCENTER | DT_SINGLELINE);

    // SelectObject(hdc, old_brush);
    SelectObject(hdc, old_pen);
    // DeleteObject(brush);
    DeleteObject(pen);
    ReleaseDC(hwnd, hdc);
}
