pub const point = @import("point.zig");
pub const Point = point.Point;
pub const input = @import("input.zig");
pub const grid = @import("grid.zig");
pub const Grid = grid.Grid;

pub fn abs(x: anytype) @TypeOf(x) {
    return if (x < 0) -x else x;
}

pub fn mod(a: anytype, b: anytype) @TypeOf(a) {
    const r = @mod(a, b);
    return if (r < 0) r + b else r;
}
