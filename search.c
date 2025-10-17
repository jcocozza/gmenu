#include <string.h>
#include <ctype.h>
#include <stdlib.h>

#include "search.h"

// search methods

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

search_results_t *search(search_config_t search_config, item_list_t *list, char *term) {
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
    if (match(list->items[i].alias, term, search_config.ignore_case)) {
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
