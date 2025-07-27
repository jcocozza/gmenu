#import <Cocoa/Cocoa.h>

void setWindowCollectionBehavior(void *window) {
  NSWindow *nsWindow = (__bridge NSWindow *)window;
  [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];

  NSUInteger behavior = NSWindowCollectionBehaviorCanJoinAllSpaces | NSWindowCollectionBehaviorFullScreenAuxiliary;
  [nsWindow setCollectionBehavior:behavior];
  //[nsWindow setLevel:CGWindowLevelForKey(kCGMaximumWindowLevelKey)-1];
  //[nsWindow orderFront:nil];
  [NSApp activateIgnoringOtherApps:YES];
  [nsWindow makeKeyAndOrderFront:nil];
}
