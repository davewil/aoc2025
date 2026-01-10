const std = @import("std");
const utils = @import("utils");

const max_neighbours_to_remove = 3;
const target_char = '@';

const directions = [_]utils.Point{
    .{ .x = 0, .y = -1 }, // Up
    .{ .x = 1, .y = -1 }, // Up-Right
    .{ .x = 1, .y = 0 }, // Right
    .{ .x = 1, .y = 1 }, // Down-Right
    .{ .x = 0, .y = 1 }, // Down
    .{ .x = -1, .y = 1 }, // Down-Left
    .{ .x = -1, .y = 0 }, // Left
    .{ .x = -1, .y = -1 }, // Up-Left
};

fn countNeighbours(g: utils.Grid(u8), x: usize, y: usize, target: u8) usize {
    var count: usize = 0;
    for (directions) |dir| {
        const nx = @as(i32, @intCast(x)) + dir.x;
        const ny = @as(i32, @intCast(y)) + dir.y;
        if (g.inBounds(nx, ny)) {
            if (g.get(@intCast(nx), @intCast(ny)) == target) {
                count += 1;
            }
        }
    }
    return count;
}

fn part1(allocator: std.mem.Allocator, g: utils.Grid(u8)) !usize {
    var removed = std.ArrayListUnmanaged(utils.Point){};
    defer removed.deinit(allocator);

    for (0..g.rows) |y| {
        for (0..g.cols) |x| {
            if (g.get(x, y) == target_char) {
                if (countNeighbours(g, x, y, target_char) <= max_neighbours_to_remove) {
                    try removed.append(allocator, .{ .x = @intCast(x), .y = @intCast(y) });
                }
            }
        }
    }
    return removed.items.len;
}

fn part2(allocator: std.mem.Allocator, g_orig: utils.Grid(u8)) !usize {
    // Clone grid because part 2 mutates it
    var g = try utils.Grid(u8).init(allocator, g_orig.rows, g_orig.cols);
    defer g.deinit();
    @memcpy(g.items, g_orig.items);

    var total_removed: usize = 0;
    var removed_in_round = std.ArrayListUnmanaged(utils.Point){};
    defer removed_in_round.deinit(allocator);

    while (true) {
        removed_in_round.clearRetainingCapacity();
        for (0..g.rows) |y| {
            for (0..g.cols) |x| {
                if (g.get(x, y) == target_char) {
                    if (countNeighbours(g, x, y, target_char) <= max_neighbours_to_remove) {
                        try removed_in_round.append(allocator, .{ .x = @intCast(x), .y = @intCast(y) });
                    }
                }
            }
        }

        if (removed_in_round.items.len == 0) break;

        for (removed_in_round.items) |p| {
            g.set(@intCast(p.x), @intCast(p.y), '.');
        }
        total_removed += removed_in_round.items.len;
    }
    return total_removed;
}

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    const input_raw = try utils.input.getPuzzleInput(allocator, 4);
    defer allocator.free(input_raw);

    const g = try utils.Grid(u8).fromLines(allocator, input_raw);
    defer {
        var mutable_g = g;
        mutable_g.deinit();
    }

    const p1 = try part1(allocator, g);
    std.debug.print("Part 1: {d}\n", .{p1});

    const p2 = try part2(allocator, g);
    std.debug.print("Part 2: {d}\n", .{p2});
}

test "example" {
    const allocator = std.testing.allocator;
    const input =
        \\..@@.@@@@.
        \\@@@.@.@.@@
        \\@@@@@.@.@@
        \\@.@@@@..@.
        \\@@.@@@@.@@
        \\.@@@@@@@.@
        \\.@.@.@.@@@
        \\@.@@@.@@@@
        \\.@@@@@@@@.
        \\@.@.@@@.@.
    ;
    const g = try utils.Grid(u8).fromLines(allocator, input);
    defer {
        var mutable_g = g;
        mutable_g.deinit();
    }

    try std.testing.expectEqual(@as(usize, 13), try part1(allocator, g));
    try std.testing.expectEqual(@as(usize, 43), try part2(allocator, g));
}
