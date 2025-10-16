#include <windows.h>
#include <stdio.h>
#include "platform.h"

static HWND hwnd = NULL;
static HINSTANCE hInstance = NULL;
static int done = 0;

static gmenu_keypress_t last_key = {KEY_NONE, 0};

int main(int argc, char **argv);

int WINAPI wWinMain(HINSTANCE hInstance, HINSTANCE hPrevInstance, LPWSTR lpCmdLine, int nCmdShow) {
    // Attach to the parent console, if this was launched from cmd/powershell
    if (AttachConsole(ATTACH_PARENT_PROCESS)) {
        FILE *fp;
        freopen_s(&fp, "CONOUT$", "w", stdout);
        freopen_s(&fp, "CONOUT$", "w", stderr);
    }

    // Convert lpCmdLine to argc/argv if needed
    int argc = __argc;
    char **argv = __argv;

    // Call your normal main function
    int result = main(argc, argv);

    return result;
}

// Forward declarations
LRESULT CALLBACK WndProc(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam);

void init() {
    hInstance = GetModuleHandle(NULL);

    const wchar_t CLASS_NAME[] = L"MinimalWindowClass";

    WNDCLASS wc = {0};
    wc.lpfnWndProc = WndProc;
    wc.hInstance = hInstance;
    wc.lpszClassName = CLASS_NAME;
    wc.hCursor = LoadCursor(NULL, IDC_ARROW);
    wc.hbrBackground = (HBRUSH)(COLOR_WINDOW+1);

    RegisterClass(&wc);

    hwnd = CreateWindowEx(
        WS_EX_TOPMOST,
        CLASS_NAME,
        L"GMenu",
        WS_OVERLAPPEDWINDOW,
        CW_USEDEFAULT, CW_USEDEFAULT, 800, 400,
        NULL, NULL, hInstance, NULL);

    if (!hwnd) {
        fprintf(stderr, "Failed to create window\n");
        return;
    }

    ShowWindow(hwnd, SW_SHOW);
    UpdateWindow(hwnd);
}

void teardown() {
    if (hwnd) {
        DestroyWindow(hwnd);
        hwnd = NULL;
    }
    done = 1;
}

int should_close() {
    MSG msg;
    last_key.k = KEY_NONE;
    last_key.c = 0;

    while (PeekMessage(&msg, NULL, 0, 0, PM_REMOVE)) {
        TranslateMessage(&msg);
        DispatchMessage(&msg);
    }

    return done;
}

gmenu_keypress_t get_key_press() {
    gmenu_keypress_t k = last_key;
    last_key.k = KEY_NONE;
    last_key.c = 0;
    return k;
}

// Draw a simple text UI
void draw(char *user_prompt, char *user_input, search_results_t *results, int result_offset, int selected_result) {
    if (!hwnd) return;

    HDC hdc = GetDC(hwnd);
    RECT rect;
    GetClientRect(hwnd, &rect);

    // Clear background
    HBRUSH brush = CreateSolidBrush(RGB(255, 255, 255));
    FillRect(hdc, &rect, brush);
    DeleteObject(brush);

    // Draw user prompt and input
    SetTextColor(hdc, RGB(0,0,0));
    SetBkMode(hdc, TRANSPARENT);
    TextOutA(hdc, 10, 10, user_prompt, (int)strlen(user_prompt));
    TextOutA(hdc, 10, 30, user_input, (int)strlen(user_input));

    // Draw search results
    if (results && results->cnt> 0) {
        for (int i = 0; i < results->cnt; i++) {
            int y = 60 + i * 20;
            if (i == selected_result) {
                HBRUSH selBrush = CreateSolidBrush(RGB(200, 200, 255));
                RECT selRect = {0, y, rect.right, y+20};
                FillRect(hdc, &selRect, selBrush);
                DeleteObject(selBrush);
            }
            TextOutA(hdc, 10, y, results->matches[i]->value, (int)strlen(results->matches[i]->value));
        }
    }

    ReleaseDC(hwnd, hdc);
}

// Window procedure
LRESULT CALLBACK WndProc(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam) {
    switch(msg) {
        case WM_CLOSE:
            done = 1;
            PostQuitMessage(0);
            return 0;

        case WM_KEYDOWN:
            switch(wParam) {
                case VK_LEFT:  last_key.k = KEY_LEFT; break;
                case VK_RIGHT: last_key.k = KEY_RIGHT; break;
                case VK_RETURN:last_key.k = KEY_ENTER; break;
                case VK_BACK:  last_key.k = KEY_BACKSPACE; break;
                default: last_key.k = KEY_OTHER; break;
            }
            return 0;

        case WM_CHAR:
            last_key.k = KEY_CHAR;
            last_key.c = (char)wParam;
            return 0;

        case WM_PAINT: {
            PAINTSTRUCT ps;
            HDC hdc = BeginPaint(hwnd, &ps);
            EndPaint(hwnd, &ps);
            return 0;
        }

        case WM_DESTROY:
            done = 1;
            PostQuitMessage(0);
            return 0;

        default:
            return DefWindowProc(hwnd, msg, wParam, lParam);
    }
}
