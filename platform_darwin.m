#include "platform.h"
#import <Cocoa/Cocoa.h>
#include <stdio.h>

static int height = 30;
static NSWindow *window = nil;
static NSView *view = nil;
static gmenu_keypress_t last_kp = {0};
static int close_flag = 0;
static NSDictionary *text_attrs = nil;
static NSImage *backbuffer = nil;

@interface OverlayWindow : NSWindow
@end

@implementation OverlayWindow
- (BOOL)canBecomeKeyWindow {
  return YES;
}
@end

@interface OverlayView : NSView
@end

@implementation OverlayView

- (BOOL)acceptsFirstResponder {
  return YES;
}

- (void)drawRect:(NSRect)dirtyRect {
  [super drawRect:dirtyRect];
  if (backbuffer) {
    [backbuffer drawInRect:[self bounds]
                  fromRect:NSZeroRect
                 operation:NSCompositingOperationCopy
                  fraction:1.0];
  }
}

- (void)keyDown:(NSEvent *)event {
  unsigned short keyCode = [event keyCode];
  last_kp.k = KEY_NONE;
  last_kp.c = 0;

  switch (keyCode) {
  case 123:
    last_kp.k = KEY_LEFT;
    break;
  case 124:
    last_kp.k = KEY_RIGHT;
    break;
  case 36:
    last_kp.k = KEY_ENTER;
    break;
  case 51:
    last_kp.k = KEY_BACKSPACE;
    break;
  case 53:
    last_kp.k = KEY_ESC;
    break;
  default:
    if ([[event characters] length] == 1) {
      unichar c = [[event characters] characterAtIndex:0];
      last_kp.c = c;
      last_kp.k = KEY_CHAR;
    }
    break;
  }
}

@end

void init() {
  @autoreleasepool {
    [NSApplication sharedApplication];
    [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];

    NSScreen *screen = [NSScreen mainScreen];
    NSRect screenRect = [screen frame];

    NSRect windowRect = NSMakeRect(0, screenRect.size.height - height,
                                   screenRect.size.width, height);

    window =
        [[OverlayWindow alloc] initWithContentRect:windowRect
                                         styleMask:NSWindowStyleMaskBorderless
                                           backing:NSBackingStoreBuffered
                                             defer:NO];

    [window setLevel:NSFloatingWindowLevel];
    [window setOpaque:NO];
    [window setBackgroundColor:[NSColor whiteColor]];
    [window setIgnoresMouseEvents:NO];
    [window makeKeyAndOrderFront:nil];
    [window setCollectionBehavior:NSWindowCollectionBehaviorCanJoinAllSpaces |
                                  NSWindowCollectionBehaviorStationary];

    view = [[OverlayView alloc] initWithFrame:windowRect];
    [window setContentView:view];
    [window makeFirstResponder:view];

    // Set up text attributes for drawing
    NSFont *font = [NSFont systemFontOfSize:13.0];
    text_attrs = @{
      NSFontAttributeName : font,
      NSForegroundColorAttributeName : [NSColor blackColor]
    };
    [text_attrs retain];

    // ensure the app is finished launching so the runloop is ready
    [NSApp finishLaunching];
    // ensure we are key (sometimes needed)
    [window makeKeyWindow];
    [NSApp activateIgnoringOtherApps:YES];
  }
}

void teardown() {
  @autoreleasepool {
    if (backbuffer) {
      [backbuffer release];
      backbuffer = nil;
    }
    if (text_attrs) {
      [text_attrs release];
      text_attrs = nil;
    }
    if (window) {
      [window close];
      [window release];
      window = nil;
    }
    if (view) {
      [view release];
      view = nil;
    }
  }
}

int should_close() { return close_flag; }

gmenu_keypress_t get_key_press() {
  @autoreleasepool {
    NSEvent *event;
    while ((event = [NSApp nextEventMatchingMask:NSEventMaskAny
                                       untilDate:[NSDate distantPast]
                                          inMode:NSDefaultRunLoopMode
                                         dequeue:YES])) {
      [NSApp sendEvent:event];
    }

    gmenu_keypress_t kp = last_kp;
    last_kp.k = KEY_NONE;
    last_kp.c = 0;
    return kp;
  }
}

int text_width(char *txt) {
  @autoreleasepool {
    NSString *str = [NSString stringWithUTF8String:txt];
    NSSize size = [str sizeWithAttributes:text_attrs];
    return (int)size.width;
  }
}

int screen_width() {
  @autoreleasepool {
    NSScreen *screen = [NSScreen mainScreen];
    NSRect screenRect = [screen frame];
    return (int)screenRect.size.width;
  }
}

void begin_draw() {
  @autoreleasepool {
    if (!backbuffer) {
      NSRect bounds = [view bounds];
      backbuffer = [[NSImage alloc] initWithSize:bounds.size];
    }
    [backbuffer lockFocus];
  }
}

void end_draw() {
  @autoreleasepool {
    [backbuffer unlockFocus];
    [view setNeedsDisplay:YES];
  }
}

void draw_text(char *txt, int x, int y, gmenu_color_t c) {
  @autoreleasepool {
    if (txt == NULL)
      return;
    NSString *str = [NSString stringWithUTF8String:txt];
    str = [str stringByReplacingOccurrencesOfString:@"\n" withString:@""];
    str = [str stringByReplacingOccurrencesOfString:@"\r" withString:@""];

    NSColor *color =
        (c == GMENU_RED) ? [NSColor redColor] : [NSColor blackColor];
    NSFont *font = [text_attrs objectForKey:NSFontAttributeName];
    NSDictionary *attrs =
        @{NSFontAttributeName : font, NSForegroundColorAttributeName : color};

    // printf("drawing %s at %d,%d\n", txt, x, y);
    NSView *cv = [window contentView];
    NSRect bounds = [cv bounds];

    // If view is not flipped, convert y from top-origin to AppKit
    // bottom-origin, and subtract font ascender so y indicates top-of-text
    CGFloat drawX = (CGFloat)x;
    CGFloat drawY = (CGFloat)y;
    if (![cv isFlipped]) {
      // 'y' is top-based; convert to bottom-based and adjust for baseline:
      drawY = bounds.size.height - (CGFloat)y - [font ascender];
    } else {
      // flipped (top-left origin): we can use y directly, but still adjust for
      // ascender
      drawY = (CGFloat)y - [font ascender];
    }

    [str drawAtPoint:NSMakePoint(drawX, drawY) withAttributes:attrs];
  }
}

void clear_screen() {
  @autoreleasepool {
    NSRect bounds = [view bounds];
    [[NSColor whiteColor] setFill];
    NSRectFill(bounds);
  }
}
