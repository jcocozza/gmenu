#pragma once

#include <stddef.h>

typedef struct item {
  char *alias;
  char *value;
} item_t;

typedef struct item_list {
  item_t *items;
  size_t cnt;
  size_t cap;
} item_list_t;

typedef struct search_results {
  item_t **matches;
  size_t cnt;
} search_results_t;

typedef struct search_config {
  int ignore_case;
} search_config_t;

search_results_t *search(search_config_t search_config, item_list_t *list, char *term);

void free_results(search_results_t *results);
