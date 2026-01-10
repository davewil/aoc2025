const std = @import("std");
const root = @import("root.zig");

pub const Point = struct {
    x: i32,
    y: i32,

    pub fn init(x: i32, y: i32) Point {
        return .{ .x = x, .y = y };
    }

    pub fn manhattanDistance(a: Point, b: Point) u32 {
        return @intCast(root.abs(a.x - b.x) + root.abs(a.y - b.y));
    }
};
