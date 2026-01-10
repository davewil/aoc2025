const std = @import("std");
const utils = @import("utils");

fn parseInput(allocator: std.mem.Allocator, input: []const u8) ![][]u8 {
    var banks = std.ArrayListUnmanaged([]u8){};
    errdefer {
        for (banks.items) |bank| allocator.free(bank);
        banks.deinit(allocator);
    }

    var it = std.mem.tokenizeAny(u8, input, " \n\r\t");
    while (it.next()) |line| {
        const nums = try allocator.alloc(u8, line.len);
        for (line, 0..) |c, i| {
            nums[i] = c - '0';
        }
        try banks.append(allocator, nums);
    }
    return banks.toOwnedSlice(allocator);
}

fn pow10(n: usize) i64 {
    var result: i64 = 1;
    var i: usize = 0;
    while (i < n) : (i += 1) {
        result *= 10;
    }
    return result;
}

fn maxJoltage(bank: []const u8, digits: usize) i64 {
    if (digits == 1) {
        var max: u8 = 0;
        for (bank) |val| {
            if (val > max) max = val;
        }
        return @intCast(max);
    }

    const search_range = bank[0 .. bank.len - (digits - 1)];
    var max_digit: u8 = 0;
    var max_idx: usize = 0;
    for (search_range, 0..) |val, i| {
        if (val > max_digit) {
            max_digit = val;
            max_idx = i;
        }
    }

    // Go's slices.Index returns the first occurrence of maxDigit in the WHOLE bank.
    // In Go: idx := slices.Index(bank, maxDigit)
    // searchRange is bank[:len(bank)-(digits-1)]
    // If maxDigit occurs before searchRange ends, slices.Index(bank, maxDigit) will find it.
    // Let's be precise.
    
    // Actually, slices.Max(searchRange) finds the value, then slices.Index(bank, maxDigit) 
    // finds the FIRST index in the whole bank with that value.
    for (bank, 0..) |val, i| {
        if (val == max_digit) {
            max_idx = i;
            break;
        }
    }

    return @as(i64, max_digit) * pow10(digits - 1) + maxJoltage(bank[max_idx + 1 ..], digits - 1);
}

fn part1(banks: [][]u8) i64 {
    var total: i64 = 0;
    for (banks) |bank| {
        total += maxJoltage(bank, 2);
    }
    return total;
}

fn part2(banks: [][]u8, target: usize) i64 {
    var total: i64 = 0;
    for (banks) |bank| {
        total += maxJoltage(bank, target);
    }
    return total;
}

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    const input_raw = try utils.input.getPuzzleInput(allocator, 3);
    defer allocator.free(input_raw);

    const input = try parseInput(allocator, input_raw);
    defer {
        for (input) |bank| allocator.free(bank);
        allocator.free(input);
    }

    const p1 = part1(input);
    std.debug.print("Part 1: {d}\n", .{p1});

    const p2 = part2(input, 12);
    std.debug.print("Part 2: {d}\n", .{p2});
}

test "example" {
    const allocator = std.testing.allocator;
    const input_str =
        \\987654321111111
        \\811111111111119
        \\234234234234278
        \\818181911112111
    ;
    const input = try parseInput(allocator, input_str);
    defer {
        for (input) |bank| allocator.free(bank);
        allocator.free(input);
    }

    try std.testing.expectEqual(@as(i64, 357), part1(input));
    try std.testing.expectEqual(@as(i64, 3121910778619), part2(input, 12));
}

