const std = @import("std");
const utils = @import("utils");

fn parseInput(raw: []const u8) !void {
    var it = std.mem.tokenizeAny(u8, raw, "\n\r");
    var total: usize = 0;
    while (it.next()) |line| {
        if (std.mem.indexOfScalar(u8, line, 'x') != null) {
            var parts = std.mem.splitScalar(u8, line, ':');
            const size_str = parts.next() orelse continue;
            const pieces_str = parts.next() orelse continue;

            var size_parts = std.mem.splitScalar(u8, size_str, 'x');
            const w_str = size_parts.next() orelse continue;
            const h_str = size_parts.next() orelse continue;

            const width = try std.fmt.parseInt(i32, std.mem.trim(u8, w_str, " "), 10);
            const height = try std.fmt.parseInt(i32, std.mem.trim(u8, h_str, " "), 10);

            var piece_it = std.mem.tokenizeAny(u8, pieces_str, " ");
            var all_parts_area: i32 = 0;
            while (piece_it.next()) |p_str| {
                const p = try std.fmt.parseInt(i32, p_str, 10);
                all_parts_area += p * 7;
            }

            if (all_parts_area <= width * height) {
                total += 1;
            }
        }
    }
    std.debug.print("Total fitting grids: {d}\n", .{total});
}

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    const raw = try utils.input.getPuzzleInput(allocator, 12);
    defer allocator.free(raw);

    try parseInput(raw);
}

test "example" {
    try parseInput("");
}
