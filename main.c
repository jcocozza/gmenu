#include "raylib.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

#ifdef __APPLE__
void setWindowBehavior(void *window);
#endif

#define MAX_LINE_LENGTH 1024
#define MAX_INPUT_CHARS 1024

struct elm {
  char *alias;
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

char *strcombine(char *s1, char *s2) {
  size_t l1 = strlen(s1);
  size_t l2 = strlen(s2);

  char *res = malloc(l1 + l2 + 1);
  if (res == NULL) {
    return NULL;
  }

  memcpy(res, s1, l1);
  memcpy(res + l1, s2, l2);
  res[l1 + l2] = '\0';
  return res;
}

// TODO: make this a dynamic allocation for each line
int readlines(FILE *f, int alias_mode) {
  char line[MAX_LINE_LENGTH];
  int alloclines = 0;
  int i = 0;

  while (fgets(line, sizeof(line), f)) {
    if (i + 1 > alloclines) {
      alloclines += 256;
      elms = realloc(elms, alloclines * sizeof(*elms));
    }

    if (alias_mode) {
      char *tok = strtok(line, " ");
      char *first = strcp(tok);
      tok = strtok(NULL, " ");
      char *second = strcp(tok);
      if (tok == NULL) {
        second = first;
      }
      elms[i].alias = first;
      elms[i].text = second;
    } else {
      char *ln = strcp(line);
      elms[i].alias = ln;
      elms[i].text = ln;
    }
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
    if (match(term, elms[i].alias)) {
      results[j].text = elms[i].text;
      results[j].alias = elms[i].alias;
      j++;
    }
  };
  *numresults = j;
  return results;
}

const int FONT_SIZE = 10;

void draw(int max_width, char *prompt, char *input, struct elm *results,
          int num_results, int result_offset, int selected_result) {

  int max_input_size = max_width;
  int min_input_size = .25 * max_input_size;

  char *final_prompt = "";
  if (strlen(prompt)) {
    final_prompt = strcombine(prompt, ": ");
  }
  char *left = strcombine(final_prompt, input);
  int left_width = MeasureText(left, FONT_SIZE);
  if (left == NULL) {
    left = "";
    left_width = 0;
  }
  DrawText(left, 0, 10, FONT_SIZE, BLACK);

  int results_size = max_width - left_width;
  int rendered_results_size = 0;
  int i = result_offset;

  int offset = min_input_size;
  if (left_width > min_input_size) {
    offset = left_width;
  }
  printf("prompt: %s, input: %s, LEFT: %s, LEFT WIDTH: %d, OFFSET: %d, RESULT "
         "SIZE: %d\n",
         prompt, input, left, left_width, offset, results_size);
  while (rendered_results_size <= results_size && i < num_results) {
    struct elm itm = results[i];

    char *display_text = itm.alias;
    int display_text_size = MeasureText(display_text, FONT_SIZE);
    rendered_results_size += display_text_size;

    if (i == selected_result) {
      DrawText(display_text, offset, 10, FONT_SIZE, RED);
    } else {
      DrawText(display_text, offset, 10, FONT_SIZE, BLACK);
    }
    i++;
    offset += display_text_size;
  }
}

void usage() { fprintf(stderr, "usage: gmenu [flags] [FILE]\n"); }

int main(int argc, char *argv[]) {
  SetTraceLogLevel(LOG_NONE); // tell raylib to be quiet

  FILE *f = stdin;

  // flags
  int alias_mode = 0;
  char *prompt = "";

  // CLI stuff
  for (int i = 1; i < argc; i++) {
    if (argv[i][0] == '-' && strlen(argv[i]) > 0) {
      if (!strcmp(argv[i], "--help")) {
        printf("help menu\n");
        return 0;
      } else if (!strcmp(argv[i], "-h")) {
        printf("help menu\n");
        return 0;
      } else if (!strcmp(argv[i], "-a") || !strcmp(argv[i], "--alias")) {
        alias_mode = 1;
      } else if (!strcmp(argv[i], "-p") || !strcmp(argv[i], "--prompt")) {
        i++;
        prompt = argv[i];
      } else { // undefined flags
        usage();
        return 1;
      }
    } else {
      f = fopen(argv[i], "r");
      if (!f) {
        perror("failed to open file");
        return 1;
      }
      break; // for now we just allow for 1 file
    }
  }

  numlines = readlines(f, alias_mode); // populate elms
  int numresults = numlines;
  struct elm *results = elms;
  int do_search = 0;

  SetConfigFlags(FLAG_WINDOW_UNDECORATED | FLAG_WINDOW_TOPMOST);
  // Create temporary window to initialize the windowing system
  InitWindow(100, 100, "Init");
#ifdef __APPLE__
  void *window = GetWindowHandle();
  setWindowBehavior(window);
#endif

  int monitor = GetCurrentMonitor();
  int width = GetMonitorWidth(monitor);
  const int maxWidth = width - 20;
  int height = 30;
  SetWindowSize(width, height);
  SetWindowPosition(0, 0);

  int inputCnt = 0;
  char input[MAX_INPUT_CHARS + 1] = "\0";

  int result_offset = 0;
  int selected_result = 0;

  // initial draw with all all elements
  draw(maxWidth, prompt, input, results, numresults, result_offset,
       selected_result);

  while (!WindowShouldClose()) {
    int key = GetCharPressed();
    while (key > 0) {
      if ((key >= 32) && (key <= 125) && (inputCnt < MAX_INPUT_CHARS)) {
        do_search = 1;
        input[inputCnt] = (char)key;
        input[inputCnt+1] = '\0';
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
    draw(maxWidth, prompt, input, results, numresults, result_offset,
         selected_result);
    EndDrawing();
  }
  CloseWindow();
  return 0;
}
