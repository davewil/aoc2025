const std = @import("std");
const utils = @import("utils");

fn part1(g_orig: utils.Grid(u8)) !usize {
    const allocator = g_orig.allocator;
    var g = try utils.Grid(u8).init(allocator, g_orig.rows, g_orig.cols);
    defer g.deinit();
    @memcpy(g.items, g_orig.items);

    var start_x: usize = 0;
    var found = false;
    for (0..g.cols) |x| {
        if (g.get(x, 0) == 'S') {
            start_x = x;
            found = true;
            break;
        }
    }
    if (!found or g.rows <= 1) return 0;

    g.set(start_x, 1, '|');
    for (1..g.rows) |y| {
        for (0..g.cols) |x| {
            if (g.get(x, y) == '|') {
                if (y + 1 < g.rows) {
                    const below = g.get(x, y + 1);
                    if (below == '^') {
                        if (x > 0) g.set(x - 1, y + 1, '|');
                        if (x + 1 < g.cols) g.set(x + 1, y + 1, '|');
                    } else if (below == '.') {
                        g.set(x, y + 1, '|');
                    }
                }
            }
        }
    }

    var collected: usize = 0;
    for (1..g.rows) |y| {
        for (0..g.cols) |x| {
            if (g.get(x, y) == '^' and g.get(x, y - 1) == '|') {
                collected += 1;
            }
        }
    }
    return collected;
}

fn part2(allocator: std.mem.Allocator, g: utils.Grid(u8)) !usize {
    var current = try allocator.alloc(usize, g.cols);
    defer allocator.free(current);
    @memset(current, 0);
    var next = try allocator.alloc(usize, g.cols);
    defer allocator.free(next);
    @memset(next, 0);

    var start_x: usize = 0;
    for (0..g.cols) |x| {
        if (g.get(x, 0) == 'S') {
            start_x = x;
            break;
        }
    }

    current[start_x] = 1;
    var range_start = start_x;
    var range_end = start_x;

    for (1..g.rows) |y| {
        @memset(next, 0);
        var new_start: usize = g.cols;
        var new_end: usize = 0;
        var any_routes = false;

        for (range_start..range_end + 1) |x| {
            if (current[x] == 0) continue;

            const cell = g.get(x, y);
            if (cell == '^') {
                if (x > 0) {
                    next[x - 1] += current[x];
                    new_start = @min(new_start, x - 1);
                    new_end = @max(new_end, x - 1);
                    any_routes = true;
                }
                if (x + 1 < g.cols) {
                    next[x + 1] += current[x];
                    new_start = @min(new_start, x + 1);
                    new_end = @max(new_end, x + 1);
                    any_routes = true;
                }
            } else if (cell == '.' or cell == 'S') {
                next[x] += current[x];
                new_start = @min(new_start, x);
                new_end = @max(new_end, x);
                any_routes = true;
            }
        }

        if (any_routes) {
            range_start = new_start;
            range_end = new_end;
        } else {
            // No routes left
            return 0;
        }

        const temp = current;
        current = next;
        next = temp;
    }

    var total: usize = 0;
    for (current) |v| total += v;
    return total;
}

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    const raw = try utils.input.getPuzzleInput(allocator, 7);
    defer allocator.free(raw);

    var g = try utils.Grid(u8).fromLines(allocator, raw);
    defer g.deinit();

    const p1 = try part1(g);
    std.debug.print("Part 1: {d}\n", .{p1});

    const p2 = try part2(allocator, g);
    std.debug.print("Part 2: {d}\n", .{p2});
}

test "example" {
    const allocator = std.testing.allocator;
    const raw =
        \\.......S.......
        \\...............
        \\.......^.......
        \\...............
        \\......^.^......
        \\...............
        \\.....^.^.^.....
        \\...............
        \\....^.^...^....
        \\...............
        \\...^.^...^.^...
        \\...............
        \\..^...^.....^..
        \\...............
        \\.^.^.^.^.^...^.
        \\...............
    ;
    var g = try utils.Grid(u8).fromLines(allocator, raw);
    defer g.deinit();

    try std.testing.expectEqual(@as(usize, 21), try part1(g));
    try std.testing.expectEqual(@as(usize, 40), try part2(allocator, g));
}
