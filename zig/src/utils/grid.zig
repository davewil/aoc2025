const std = @import("std");

pub fn Grid(comptime T: type) type {
    return struct {
        const Self = @This();
        items: []T,
        rows: usize,
        cols: usize,
        allocator: std.mem.Allocator,

        pub fn init(allocator: std.mem.Allocator, rows: usize, cols: usize) !Self {
            const items = try allocator.alloc(T, rows * cols);
            return Self{
                .items = items,
                .rows = rows,
                .cols = cols,
                .allocator = allocator,
            };
        }

        pub fn deinit(self: *Self) void {
            self.allocator.free(self.items);
        }

        pub fn get(self: Self, x: usize, y: usize) T {
            return self.items[y * self.cols + x];
        }

        pub fn set(self: *Self, x: usize, y: usize, val: T) void {
            self.items[y * self.cols + x] = val;
        }

        pub fn inBounds(self: Self, x: i32, y: i32) bool {
            return x >= 0 and x < @as(i32, @intCast(self.cols)) and y >= 0 and y < @as(i32, @intCast(self.rows));
        }

        pub fn fromLines(allocator: std.mem.Allocator, input: []const u8) !Self {
            var lines = std.mem.tokenizeAny(u8, input, "\n\r");
            var row_list = std.ArrayListUnmanaged([]const u8){};
            defer row_list.deinit(allocator);

            while (lines.next()) |line| {
                try row_list.append(allocator, line);
            }

            if (row_list.items.len == 0) return error.EmptyGrid;
            const rows = row_list.items.len;
            const cols = row_list.items[0].len;

            var grid = try Self.init(allocator, rows, cols);
            for (row_list.items, 0..) |line, y| {
                const len = @min(line.len, cols);
                for (line[0..len], 0..) |char, x| {
                    if (T == u8) {
                        grid.set(x, y, char);
                    } else if (T == i32) {
                        grid.set(x, y, @intCast(char - '0'));
                    } else {
                        @compileError("Unsupported Grid type");
                    }
                }
                // Pad with space if line is shorter than cols
                if (line.len < cols) {
                    for (line.len..cols) |x| {
                        if (T == u8) {
                            grid.set(x, y, ' ');
                        } else if (T == i32) {
                            grid.set(x, y, 0);
                        }
                    }
                }
            }
            return grid;
        }

        pub fn row(self: Self, y: usize) []T {
            return self.items[y * self.cols .. (y + 1) * self.cols];
        }

        pub fn rotateLeft(self: Self) !Self {
            var new_grid = try Self.init(self.allocator, self.cols, self.rows);
            for (0..self.rows) |y| {
                for (0..self.cols) |x| {
                    // (x, y) in old grid
                    // Rotate left: new_x = y, new_y = (cols - 1) - x
                    new_grid.set(y, self.cols - 1 - x, self.get(x, y));
                }
            }
            return new_grid;
        }
    };
}
