// this is the "api" for interacting with different platforms
#pragma once

#include "search.h"

typedef enum key {
  KEY_NONE,
  KEY_OTHER,
  KEY_CHAR,
  KEY_BACKSPACE,
  KEY_LEFT,
  KEY_RIGHT,
  KEY_ENTER,
  KEY_ESC
} gmenu_key_t;

typedef struct keypress {
  gmenu_key_t k;
  char c;
} gmenu_keypress_t;

// start up
void init();
// tear down
void teardown();
// return 1 when the app is done
int should_close();
// return key pressed by user (only the ones we care about)
gmenu_keypress_t get_key_press();
// width of text in pixels
int text_width(char *txt);
// width of the screen that is displaying the tool
int screen_width();
// draw text at x, y
void draw_text(char *txt, int x, int y);

// draw results
//void draw(char *user_prompt, char *user_input, search_results_t *results, int result_offset, int selected_result);
