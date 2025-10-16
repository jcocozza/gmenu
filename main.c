#include "platform.h"
#include "search.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

const int FONT_SIZE = 10;
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

int main(int argc, char *argv[]) {
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
      } else if (!strcmp(argv[i], "-i") ||
                 !strcmp(argv[i], "--no-ignore-case")) {
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

  // search config
  search_config_t sc = {.ignore_case = ignore_case};

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

  while (!should_close()) {
    // this is basically copied straight from
    // https://www.raylib.com/examples/text/loader.html?name=text_input_box
    gmenu_keypress_t kp = get_key_press();
    printf("key press: %d\n", kp.k);
    search_results_t *results = search(sc, list, input);

    switch (kp.k) {
    case KEY_NONE:
      break;
    case KEY_OTHER:
      break;
    case KEY_CHAR:
      if (strlen(input) >= input_count) {
        input = realloc(input, strlen(input) + 256);
        if (!input) {
          perror("malloc");
          exit(1);
        }
      }
      input[input_count] = kp.c;
      input[input_count + 1] = '\0';
      input_count++;
      break;
    case KEY_BACKSPACE:
      if (input_count > 0) {
        input_count--;
        input[input_count] = '\0';
      }
      break;
    case KEY_LEFT:
      if (results->cnt == 0) {
        break;
      };
      if (result_offset > 0) {
        result_offset--;
        selected_result--;
      } else if (result_offset <= 0) { // to end
        result_offset = results->cnt - 1;
        selected_result = results->cnt - 1;
      }
      break;
    case KEY_RIGHT:
      if (results->cnt == 0) {
        break;
      };
      if (result_offset < results->cnt - 1) {
        result_offset++;
        selected_result++;
      } else { // back to beginning
        result_offset = 0;
        selected_result = 0;
      }
      break;
    case KEY_ENTER:
      printf("%s", results->matches[selected_result]->value);
      return 0;
    }

    draw(prompt, input, results, result_offset, selected_result);
    free_results(results);
  }
  teardown();
  return 0;
}
