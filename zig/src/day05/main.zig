const std = @import("std");
const utils = @import("utils");

const Range = struct {
    start: i32,
    end: i32,

    pub fn len(self: Range) i32 {
        return self.end - self.start + 1;
    }
};

const Input = struct {
    ranges: []Range,
    numbers: []i32,
};

fn mergeRanges(allocator: std.mem.Allocator, ranges: []Range) ![]Range {
    if (ranges.len == 0) return try allocator.dupe(Range, ranges);

    // Sort by start
    std.mem.sort(Range, ranges, {}, struct {
        fn lessThan(_: void, a: Range, b: Range) bool {
            return a.start < b.start;
        }
    }.lessThan);

    var merged = std.ArrayListUnmanaged(Range){};
    errdefer merged.deinit(allocator);

    try merged.append(allocator, ranges[0]);

    for (ranges[1..]) |current| {
        var last = &merged.items[merged.items.len - 1];
        if (current.start <= last.end + 1) {
            if (current.end > last.end) {
                last.end = current.end;
            }
        } else {
            try merged.append(allocator, current);
        }
    }

    return merged.toOwnedSlice(allocator);
}

fn parseInput(allocator: std.mem.Allocator, input_str: []const u8) !Input {
    var range_list = std.ArrayListUnmanaged(Range){};
    var num_list = std.ArrayListUnmanaged(i32){};
    errdefer {
        range_list.deinit(allocator);
        num_list.deinit(allocator);
    }

    var lines = std.mem.tokenizeAny(u8, input_str, "\n\r");
    while (lines.next()) |line| {
        if (line.len == 0) continue;
        if (std.mem.indexOfScalar(u8, line, '-') != null) {
            var it = std.mem.splitScalar(u8, line, '-');
            const s_str = it.next().?;
            const e_str = it.next().?;
            try range_list.append(allocator, .{
                .start = try std.fmt.parseInt(i32, s_str, 10),
                .end = try std.fmt.parseInt(i32, e_str, 10),
            });
        } else {
            try num_list.append(allocator, try std.fmt.parseInt(i32, line, 10));
        }
    }

    const merged = try mergeRanges(allocator, range_list.items);
    range_list.deinit(allocator);

    return Input{
        .ranges = merged,
        .numbers = try num_list.toOwnedSlice(allocator),
    };
}

fn part1(input: Input) i32 {
    var count: i32 = 0;
    for (input.numbers) |num| {
        // Binary search for the first range that ends at or after num
        var low: usize = 0;
        var high: usize = input.ranges.len;
        while (low < high) {
            const mid = low + (high - low) / 2;
            if (input.ranges[mid].end < num) {
                low = mid + 1;
            } else {
                high = mid;
            }
        }

        if (low < input.ranges.len and input.ranges[low].start <= num) {
            count += 1;
        }
    }
    return count;
}

fn part2(input: Input) i32 {
    var count: i32 = 0;
    for (input.ranges) |rng| {
        count += rng.len();
    }
    return count;
}

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    const input_raw = try utils.input.getPuzzleInput(allocator, 5);
    defer allocator.free(input_raw);

    const input = try parseInput(allocator, input_raw);
    defer {
        allocator.free(input.ranges);
        allocator.free(input.numbers);
    }

    const p1 = part1(input);
    std.debug.print("Part 1: {d}\n", .{p1});

    const p2 = part2(input);
    std.debug.print("Part 2: {d}\n", .{p2});
}

test "example" {
    const allocator = std.testing.allocator;
    const input_str =
        \\3-5
        \\10-14
        \\16-20
        \\12-18
        \\
        \\1
        \\5
        \\8
        \\11
        \\17
        \\32
    ;
    const input = try parseInput(allocator, input_str);
    defer {
        allocator.free(input.ranges);
        allocator.free(input.numbers);
    }

    try std.testing.expectEqual(@as(usize, 2), input.ranges.len);
    try std.testing.expectEqual(@as(i32, 3), input.ranges[0].start);
    try std.testing.expectEqual(@as(i32, 5), input.ranges[0].end);
    try std.testing.expectEqual(@as(i32, 10), input.ranges[1].start);
    try std.testing.expectEqual(@as(i32, 20), input.ranges[1].end);

    try std.testing.expectEqual(@as(i32, 3), part1(input));
    try std.testing.expectEqual(@as(i32, 14), part2(input));
}
