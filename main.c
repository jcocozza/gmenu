#include "raylib.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#ifdef __APPLE__
void setWindowBehavior(void *window);
#endif

#define MAX_LINE_LENGTH 1024
#define MAX_INPUT_CHARS 1024

struct elm {
  char *text;
};

static struct elm *elms = NULL;
static int numlines = 0;

char *strcp(char *s) {
  if (s == NULL)
    return NULL;
  size_t len = strlen(s) + 1;
  char *cp = malloc(len);
  if (cp) {
    memcpy(cp, s, len);
  }
  return cp;
}

char *substr(char *s, int start, int end) {
  size_t src_len = strlen(s);

  if (start > src_len)
    return NULL;
  if (end > src_len)
    end = src_len;

  int length = end - start;
  char *sub = malloc(length);

  strncpy(sub, s + start, length);
  return sub;
}

// TODO: make this a dynamic allocation for each line
int readlines() {
  char line[MAX_LINE_LENGTH];
  int alloclines = 0;
  int i = 0;

  while (fgets(line, sizeof(line), stdin)) {
    if (i + 1 > alloclines) {
      alloclines += 256;
      elms = realloc(elms, alloclines * sizeof(*elms));
    }
    elms[i].text = strcp(line);
    i++;
  }
  return i;
}

// TODO: we can make this very fancy
int match(char *search, char *line) { return strstr(line, search) != NULL; }

struct elm *search(char *term, int *numresults) {
  struct elm *results = NULL;
  int allocresults = 0;
  int j = 0;
  for (int i = 0; i < numlines; i++) {
    if (j + 1 > allocresults) {
      allocresults += 256;
      results = realloc(results, allocresults * sizeof(*results));
    }
    if (match(term, elms[i].text)) {
      results[j].text = elms[i].text;
      j++;
    }
  };
  *numresults = j;
  return results;
}

const int FONT_SIZE = 10;

void draw(int max_width, char *input, struct elm *results, int num_results,
          int result_offset, int selected_result) {
  // TODO: this should be moved to main since it doesn't need to be recomputed
  // each time
  int sep_size = MeasureText(", ", FONT_SIZE);
  int max_input_size = max_width;
  int min_input_size = .25 * max_width;

  int input_size = MeasureText(input, FONT_SIZE);
  int results_size = max_width - input_size;

  char *display_input = input;
  if (input_size > max_input_size) {
    display_input = substr(input, 0, max_input_size);
  }

  DrawText(display_input, 10, 10, FONT_SIZE, BLACK);

  int offset = min_input_size;
  if (input_size > min_input_size) {
    offset = input_size;
  }

  int curr_results_size = 0;

  int i = result_offset;
  while (curr_results_size <= results_size && i < num_results) {
    struct elm itm = results[i];

    char *display_text = itm.text;
    if (strlen(display_text) == 0 || strcmp(display_text, "\n") == 0 ||
        strcmp(display_text, " ") == 0) {
      display_text = "_";
    }
    int txt_size = MeasureText(display_text, FONT_SIZE);
    curr_results_size += txt_size;

    // TODO: include separator
    if (i == selected_result) { // highlighted one
      DrawText(display_text, offset + sep_size, 10, FONT_SIZE, RED);
    } else { // other results
      DrawText(display_text, offset + sep_size, 10, FONT_SIZE, BLACK);
    }
    i++;

    offset += txt_size + sep_size;
  }
}

int main(void) {
  SetTraceLogLevel(LOG_NONE);
  numlines = readlines(); // populate elms
  int numresults = numlines;
  struct elm *results = elms;
  int do_search = 0;

  SetConfigFlags(FLAG_WINDOW_UNDECORATED | FLAG_WINDOW_TOPMOST);
  // Create temporary window to initialize the windowing system
  InitWindow(100, 100, "Init");
#ifdef __APPLE__
  void *window = GetWindowHandle(); // Returns NSWindow* on macOS
  setWindowBehavior(window);
#endif

  int monitor = GetCurrentMonitor();
  int width = GetMonitorWidth(monitor);
  const int maxWidth = width - 20;
  int height = 30;
  // Re-set the window size (since we already created a dummy one)
  SetWindowSize(width, height);
  // SetWindowTitle("Thin Bar");
  SetWindowPosition(0, 0);

  int inputCnt = 0;
  char input[MAX_INPUT_CHARS] = "";
  int result_offset = 0;
  int selected_result = 0;

  draw(maxWidth, input, results, numresults, result_offset, selected_result);

  while (!WindowShouldClose()) {

    int key = GetCharPressed();
    while (key > 0) {
      if ((key >= 32) && (key <= 125) && (inputCnt < MAX_INPUT_CHARS)) {
        do_search = 1;
        input[inputCnt] = (char)key;
        inputCnt++;
      }
      key = GetCharPressed();
    }
    if (IsKeyPressed(KEY_BACKSPACE)) {
      do_search = 1;
      if (inputCnt > 0) {
        inputCnt--;
        input[inputCnt] = '\0';
      }
    }

    if (do_search) {
      results = search(input, &numresults);
      do_search = 0;
      // back to start when we do a new search
      selected_result = 0;
      result_offset = 0;
    }

    if (IsKeyPressed(KEY_RIGHT) && numresults != 0) {
      if (result_offset < numresults) {
        result_offset++;
        selected_result++;
      } else if (result_offset >= numresults) { // back to beginning
        result_offset = 0;
        selected_result = 0;
      }
    }
    if (IsKeyPressed(KEY_LEFT) && numresults != 0) {
      if (result_offset > 0) {
        result_offset--;
        selected_result--;
      } else if (result_offset <= 0) { // to end
        result_offset = numresults - 1;
        selected_result = numresults - 1;
      }
    }

    if (IsKeyPressed(KEY_ENTER)) {
      printf("%s", results[selected_result].text);
      return 0;
    }

    BeginDrawing();
    ClearBackground(RAYWHITE);
    draw(maxWidth, input, results, numresults, result_offset, selected_result);
    EndDrawing();
  }
  CloseWindow();
  return 0;
}
