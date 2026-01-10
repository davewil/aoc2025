const std = @import("std");

/// Simple input reader that reads from the inputs/ directory.
pub fn getPuzzleInput(allocator: std.mem.Allocator, day: u8) ![]u8 {
    var path_buf: [32]u8 = undefined;
    const path = try std.fmt.bufPrint(&path_buf, "inputs/day{d:0>2}.txt", .{day});
    
    const file = try std.fs.cwd().openFile(path, .{});
    defer file.close();

    return try file.readToEndAlloc(allocator, 1024 * 1024);
}
