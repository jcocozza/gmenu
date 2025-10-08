#import <Cocoa/Cocoa.h>

void setWindowBehavior(void *window) {
    NSWindow *nsWindow = (__bridge NSWindow *)window;
    [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];

    NSUInteger behavior = NSWindowCollectionBehaviorCanJoinAllSpaces |
                          NSWindowCollectionBehaviorFullScreenAuxiliary;
    [nsWindow setCollectionBehavior:behavior];

    [NSApp activateIgnoringOtherApps:YES];
    [nsWindow makeKeyAndOrderFront:nil];
}
