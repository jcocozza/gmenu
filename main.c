#include <_string.h>
#define _POSIX_C_SOURCE 200809L

#include "raylib.h"
#include <ctype.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#ifdef __APPLE__
void setWindowBehavior(void *window);
#endif

// flags
int alias_mode = 0;
char *prompt = "";
char *delim = " ";
int ignore_case = 1;

FILE **files = NULL;

void usage() { fprintf(stderr, "usage: gmenu [flags] [FILES]\n"); }

void usage_long() {
  printf("usage: gmenu [flags] [FILE]\n");
  printf("\n");
  printf("flags:\n");
  printf("-p, --prompt: include a prompt in the gui (default: \"%s\")\n",
         prompt);
  printf("-a, --alias: use gmenu in alias mode. split each line on delim in "
         "key, value. will show key and return value\n");
  printf("-d, --delim: delim to split on. only affects alias mode. (default: "
         "\"%s\")\n",
         delim);
  printf("-i, --no-ignore-case: ignore case (default: true)");
}

typedef struct file_list {
  FILE **files;
  size_t cnt;
} file_list_t;

int add_file(file_list_t *file_list, FILE *f) {
  file_list->files =
      realloc(file_list->files, (file_list->cnt + 1) * sizeof(FILE *));
  if (!file_list->files) {
    return -1;
  }
  file_list->files[file_list->cnt] = f;
  file_list->cnt++;
  return 0;
}

typedef struct item {
  char *alias;
  char *value;
} item_t;

typedef struct item_list {
  item_t *items;
  size_t cnt;
  size_t cap;
} item_list_t;

item_list_t *create_items(void) {
  item_list_t *list = malloc(sizeof(item_list_t));
  if (!list) {
    return NULL;
  }
  list->cap = 256;
  list->cnt = 0;
  list->items = malloc(list->cap * sizeof(item_t));
  if (!list->items) {
    free(list);
    return NULL;
  }
  return list;
}

int add_item(item_list_t *list, char *alias, char *value) {
  if (list->cnt >= list->cap) {
    list->cap = list->cap * 2;
    list->items = realloc(list->items, sizeof(item_t) * list->cap);
    if (!list->items) {
      return -1;
    }
  }
  list->items[list->cnt].alias = strdup(alias);
  list->items[list->cnt].value = strdup(value);
  if (!list->items[list->cnt].alias || !list->items[list->cnt].value) {
    free(list->items[list->cnt].alias);
    free(list->items[list->cnt].value);
    return -1;
  }
  list->cnt++;
  return 0;
}

typedef struct str_split {
  char *first;
  char *second;
} str_split_t;

str_split_t *strsplit(char *s, char *delim) {
  str_split_t *ss = malloc(sizeof(str_split_t));
  if (!ss)
    return NULL;

  char *match = strstr(s, delim);
  if (!match) {
    ss->first = strdup(s);
    ss->second =
        strdup(""); // Or strdup(s) if you want both to be the whole string
    return ss;
  }

  int first_size = match - s; // Pointer arithmetic is cleaner
  ss->first = malloc(first_size + 1);
  if (!ss->first) {
    free(ss);
    return NULL;
  }
  strncpy(ss->first, s, first_size);
  ss->first[first_size] = '\0';

  ss->second = strdup(match + strlen(delim)); // Skip delimiter
  if (!ss->second) {
    free(ss->first);
    free(ss);
    return NULL;
  }
  return ss;
}

int MAX_LINE_LEN = 1024;
void readfile(FILE *f, item_list_t *item_list) {
  char line[MAX_LINE_LEN];

  while (fgets(line, sizeof(line), f)) {
    if (alias_mode) {
      str_split_t *split = strsplit(line, delim);
      add_item(item_list, split->first, split->second);
      free(split->first);
      free(split->second);
      free(split);
    } else {
      add_item(item_list, line, line);
    }
  }
}

typedef struct search_results {
  item_t **matches;
  size_t cnt;
} search_results_t;

// 0 false
// 1 true
int match(char *s1, char *s2, int ignore_case) {
  char *c1 = strdup(s1);
  char *c2 = strdup(s2);
  if (!c1 || !c2) {
    free(c1);
    free(c2);
    return 0;
  }
  // always match if one is empty
  if (!strcmp(s1, "") || !strcmp(s2, "")) {
    return 1;
  }
  if (ignore_case) {
    for (int i = 0; c1[i]; i++) {
      c1[i] = tolower(c1[i]);
    }
    for (int i = 0; c2[i]; i++) {
      c2[i] = tolower(c2[i]);
    }
  }
  int res = strstr(c1, c2) != NULL;
  free(c1);
  free(c2);
  return res;
}

search_results_t *search(item_list_t *list, char *term) {
  item_t **matches = malloc(list->cnt * sizeof(item_t *));
  if (!matches)
    return NULL;

  search_results_t *results = malloc(sizeof(search_results_t));
  if (!results) {
    free(matches);
    return NULL;
  }
  results->matches = matches;
  results->cnt = 0;
  for (size_t i = 0; i < list->cnt; i++) {
    if (match(list->items[i].alias, term, ignore_case)) {
      results->matches[results->cnt] = &list->items[i];
      results->cnt++;
    }
  }
  // shrink to total search results not size of entire search space
  results->matches = realloc(results->matches, sizeof(item_t *) * results->cnt);
  return results;
}

void free_results(search_results_t *results) {
  free(results->matches);
  free(results);
}

const int FONT_SIZE = 10;
void draw(int max_width, char *user_prompt, char *user_input,
          search_results_t *results, int result_offset, int selected_result) {

  // TODO: this should not be in the render loop
  int spacer_width = MeasureText("  ", FONT_SIZE);
  int max_input_size = max_width;
  int min_input_size = .25 * max_input_size;

  int prompt_len = strlen(user_prompt);
  if (prompt_len != 0) {
    prompt_len += 2; // 2 for ": "
  }
  // i am not 100% sure that this is an okay thing to do.
  // it should probably be a heap allocation
  size_t final_size = prompt_len + strlen(user_input) + 1;
  char final_prompt[final_size];
  if (prompt_len == 0) {
    snprintf(final_prompt, final_size, "%s", user_input);
  } else {
    snprintf(final_prompt, final_size, "%s: %s", user_prompt, user_input);
  }
  DrawText(final_prompt, 10, 10, FONT_SIZE, BLACK);

  int prompt_width = MeasureText(final_prompt, FONT_SIZE);
  int results_width = max_width - prompt_width;
  int rendered_results_width = 0;
  int i = result_offset;
  int offset = min_input_size;
  if (prompt_width > min_input_size) {
    offset = prompt_width;
  }
  offset += spacer_width;

  while (rendered_results_width <= results_width && i < results->cnt) {
    char *display_text = results->matches[i]->alias;
    if (isspace(display_text[0]) &&
        strlen(display_text) ==
            1) { // this is dumb. it doesn't do it's job for "  "
      display_text = "<whitespace>";
    }

    int display_text_width = MeasureText(display_text, FONT_SIZE);
    rendered_results_width += display_text_width + spacer_width;
    if (i == selected_result) {
      DrawText(display_text, offset, 10, FONT_SIZE, RED);
    } else {
      DrawText(display_text, offset, 10, FONT_SIZE, BLACK);
    }
    i++;
    offset += display_text_width + spacer_width;
  }
}

int main(int argc, char *argv[]) {
  SetTraceLogLevel(LOG_NONE); // tell raylib to be quiet
  item_list_t *list = create_items();
  file_list_t *file_list = malloc(sizeof(file_list_t *));
  if (!file_list) {
    perror("bozo");
    exit(1);
  }
  file_list->cnt = 0;
  file_list->files = NULL;

  // CLI stuff
  for (int i = 1; i < argc; i++) {
    if (argv[i][0] == '-' && strlen(argv[i]) > 0) {
      if (!strcmp(argv[i], "--help")) {
        usage_long();
        return 0;
      } else if (!strcmp(argv[i], "-h")) {
        printf("help menu\n");
        return 0;
      } else if (!strcmp(argv[i], "-a") || !strcmp(argv[i], "--alias")) {
        alias_mode = 1;
      } else if (!strcmp(argv[i], "-p") || !strcmp(argv[i], "--prompt")) {
        i++;
        prompt = argv[i];
      } else if (!strcmp(argv[i], "-d") || !strcmp(argv[i], "--delim")) {
        i++;
        delim = argv[i];
      } else if (!strcmp(argv[i], "-i") || !strcmp(argv[i], "--no-ignore-case")) {
        ignore_case = 0;
      } else { // undefined flags
        usage_long();
        return 1;
      }
    } else {
      FILE *f = fopen(argv[i], "r");
      if (!f) {
        perror("failed to open file");
        return 1;
      }
      if (add_file(file_list, f)) {
        perror("and you fail");
        exit(1);
      }
    }
  }
  if (file_list->cnt == 0) {
    if (add_file(file_list, stdin)) {
      perror("and you fail");
      exit(1);
    }
  }

  for (size_t i = 0; i < file_list->cnt; i++) {
    readfile(file_list->files[i], list);
    fclose(file_list->files[i]);
  }
  free(file_list->files);
  free(file_list);

  SetConfigFlags(FLAG_WINDOW_UNDECORATED | FLAG_WINDOW_TOPMOST);
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

  // user input
  // note to self; if we had used malloc here, would have had to manuall set
  // input[0] = '\0'
  char *input = calloc(256, 1);
  int input_count = 0;
  if (!input) {
    perror("malloc");
    exit(1);
  }
  int selected_result = 0;
  int result_offset = 0;

  while (!WindowShouldClose()) {
    // this is basically copied straight from
    // https://www.raylib.com/examples/text/loader.html?name=text_input_box
    int key = GetCharPressed();
    while (key > 0) {
      if ((key >= 32) && (key <= 125)) {
        if (strlen(input) >= input_count) {
          input = realloc(input, strlen(input) + 256);
          if (!input) {
            perror("malloc");
            exit(1);
          }
        }

        input[input_count] = (char)key;
        input[input_count + 1] = '\0';
        input_count++;
      }
      key = GetCharPressed();
    }
    if (IsKeyPressed(KEY_BACKSPACE)) {
      if (input_count > 0) {
        input_count--;
        input[input_count] = '\0';
      }
    }
    search_results_t *results = search(list, input);
    // back to start when we do a new search
    // selected_result = 0;
    // result_offset = 0;

    if (IsKeyPressed(KEY_RIGHT) && results->cnt != 0) {
      if (result_offset < results->cnt - 1) {
        result_offset++;
        selected_result++;
      } else { // back to beginning
        result_offset = 0;
        selected_result = 0;
      }
    }
    if (IsKeyPressed(KEY_LEFT) && results->cnt != 0) {
      if (result_offset > 0) {
        result_offset--;
        selected_result--;
      } else if (result_offset <= 0) { // to end
        result_offset = results->cnt - 1;
        selected_result = results->cnt - 1;
      }
    }

    if (IsKeyPressed(KEY_ENTER)) {
      printf("%s", results->matches[selected_result]->value);
      return 0;
    }

    BeginDrawing();
    ClearBackground(RAYWHITE);
    draw(maxWidth, prompt, input, results, result_offset, selected_result);
    EndDrawing();
    free_results(results);
  }
  CloseWindow();
  return 0;
}
