const std = @import("std");
const utils = @import("utils");

const Range = struct {
    start: i64,
    end: i64,
};

fn parseInput(allocator: std.mem.Allocator, input: []const u8) ![]Range {
    var ranges = std.ArrayListUnmanaged(Range){};
    errdefer ranges.deinit(allocator);

    var it = std.mem.tokenizeAny(u8, input, " \n\r\t,");
    while (it.next()) |token| {
        if (token.len == 0) continue;
        var parts = std.mem.splitScalar(u8, token, '-');
        const start_str = parts.next() orelse continue;
        const end_str = parts.next() orelse continue;
        const start = try std.fmt.parseInt(i64, start_str, 10);
        const end = try std.fmt.parseInt(i64, end_str, 10);
        try ranges.append(allocator, .{ .start = start, .end = end });
    }
    return ranges.toOwnedSlice(allocator);
}

fn part1(ranges: []const Range) i64 {
    var total: i64 = 0;
    var buf: [32]u8 = undefined;
    for (ranges) |r| {
        var i = r.start;
        while (i <= r.end) : (i += 1) {
            const s = std.fmt.bufPrint(&buf, "{d}", .{i}) catch continue;
            if (s.len % 2 != 0) continue;
            const mid = s.len / 2;
            if (std.mem.eql(u8, s[0..mid], s[mid..])) {
                total += i;
            }
        }
    }
    return total;
}

fn getDivisors(allocator: std.mem.Allocator, n: usize) ![]usize {
    var result = std.ArrayListUnmanaged(usize){};
    errdefer result.deinit(allocator);
    var i: usize = 1;
    while (i <= n / 2) : (i += 1) {
        if (n % i == 0) {
            try result.append(allocator, i);
        }
    }
    return result.toOwnedSlice(allocator);
}

fn hasRepeatingPattern(s: []const u8, chunk_len: usize) bool {
    if (chunk_len <= 0 or s.len % chunk_len != 0) return false;
    const pattern = s[0..chunk_len];
    var i = chunk_len;
    while (i < s.len) : (i += chunk_len) {
        if (!std.mem.eql(u8, s[i .. i + chunk_len], pattern)) return false;
    }
    return true;
}

fn part2(allocator: std.mem.Allocator, ranges: []const Range) !i64 {
    var seen = std.AutoHashMap(i64, void).init(allocator);
    defer seen.deinit();
    var total: i64 = 0;
    var buf: [32]u8 = undefined;

    for (ranges) |r| {
        var i = r.start;
        while (i <= r.end) : (i += 1) {
            if (seen.contains(i)) continue;

            const s = std.fmt.bufPrint(&buf, "{d}", .{i}) catch continue;
            if (s.len < 2) continue;

            const divs = try getDivisors(allocator, s.len);
            defer allocator.free(divs);

            for (divs) |d| {
                if (hasRepeatingPattern(s, d)) {
                    try seen.put(i, {});
                    total += i;
                    break;
                }
            }
        }
    }
    return total;
}

// Concurrent versions
fn part1Worker(range: Range, result: *i64) void {
    var total: i64 = 0;
    var buf: [32]u8 = undefined;
    var i = range.start;
    while (i <= range.end) : (i += 1) {
        const s = std.fmt.bufPrint(&buf, "{d}", .{i}) catch continue;
        if (s.len % 2 != 0) continue;
        const mid = s.len / 2;
        if (std.mem.eql(u8, s[0..mid], s[mid..])) {
            total += i;
        }
    }
    result.* = total;
}

fn part1Concurrent(allocator: std.mem.Allocator, ranges: []const Range) !i64 {
    var threads = try allocator.alloc(std.Thread, ranges.len);
    defer allocator.free(threads);
    var results = try allocator.alloc(i64, ranges.len);
    defer allocator.free(results);

    for (ranges, 0..) |r, i| {
        threads[i] = try std.Thread.spawn(.{}, part1Worker, .{ r, &results[i] });
    }

    for (threads) |t| t.join();

    var total: i64 = 0;
    for (results) |res| total += res;
    return total;
}

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    const input_raw = try utils.input.getPuzzleInput(allocator, 2);
    defer allocator.free(input_raw);

    const ranges = try parseInput(allocator, input_raw);
    defer allocator.free(ranges);

    const start1 = std.time.nanoTimestamp();
    const p1 = part1(ranges);
    const end1 = std.time.nanoTimestamp();
    std.debug.print("Part 1: {d} (took {d}ms)\n", .{ p1, @divTrunc(end1 - start1, 1000000) });

    const start2 = std.time.nanoTimestamp();
    const p2 = try part2(allocator, ranges);
    const end2 = std.time.nanoTimestamp();
    std.debug.print("Part 2: {d} (took {d}ms)\n", .{ p2, @divTrunc(end2 - start2, 1000000) });

    const start1c = std.time.nanoTimestamp();
    const p1c = try part1Concurrent(allocator, ranges);
    const end1c = std.time.nanoTimestamp();
    std.debug.print("Part 1 Concurrent: {d} (took {d}ms)\n", .{ p1c, @divTrunc(end1c - start1c, 1000000) });
}

test "example" {
    const allocator = std.testing.allocator;
    const input = "11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124";
    const ranges = try parseInput(allocator, input);
    defer allocator.free(ranges);

    try std.testing.expectEqual(@as(i64, 1227775554), part1(ranges));
    const p2_val = try part2(allocator, ranges);
    try std.testing.expectEqual(@as(i64, 4174379265), p2_val);
}
