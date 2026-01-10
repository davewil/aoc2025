const std = @import("std");
const utils = @import("utils");

const Instruction = struct {
    direction: u8,
    distance: i32,
};

fn crossesZero(old_mod: i32, new_mod: i32, delta: i32) bool {
    if (old_mod == 0 or delta == 0 or @mod(delta, 100) == 0) {
        return false;
    }
    if (delta > 0) {
        return new_mod <= old_mod;
    }
    return new_mod > old_mod or new_mod == 0;
}

fn parseInput(allocator: std.mem.Allocator, input: []const u8) ![]Instruction {
    var instructions = std.ArrayListUnmanaged(Instruction){};
    errdefer instructions.deinit(allocator);

    var lines = std.mem.tokenizeAny(u8, input, " \n\r\t,");
    while (lines.next()) |line| {
        if (line.len < 2) continue;
        const dir = line[0];
        const dist = try std.fmt.parseInt(i32, line[1..], 10);
        try instructions.append(allocator, .{ .direction = dir, .distance = dist });
    }
    return instructions.toOwnedSlice(allocator);
}

fn part1(instructions: []const Instruction) i32 {
    var current_notch: i32 = 50;
    var notch_zero_seen_count: i32 = 0;
    for (instructions) |instr| {
        var delta = instr.distance;
        if (instr.direction == 'L') {
            delta = -delta;
        }
        current_notch += delta;
        if (utils.mod(current_notch, 100) == 0) {
            notch_zero_seen_count += 1;
        }
    }
    return notch_zero_seen_count;
}

fn part2(instructions: []const Instruction) i32 {
    var current_notch: i32 = 50;
    var notch_zero_seen_count: i32 = 0;
    for (instructions) |instr| {
        var delta = instr.distance;
        if (instr.direction == 'L') {
            delta = -delta;
        }

        notch_zero_seen_count += @divTrunc(utils.abs(delta), 100);

        const old_mod = utils.mod(current_notch, 100);
        current_notch += delta;
        const new_mod = utils.mod(current_notch, 100);
        if (crossesZero(old_mod, new_mod, delta)) {
            notch_zero_seen_count += 1;
        }
    }
    return notch_zero_seen_count;
}

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    const input = try utils.input.getPuzzleInput(allocator, 1);
    defer allocator.free(input);

    const instructions = try parseInput(allocator, input);
    defer allocator.free(instructions);

    const p1 = part1(instructions);
    const p2 = part2(instructions);

    std.debug.print("Part 1: {d}\n", .{p1});
    std.debug.print("Part 2: {d}\n", .{p2});
}

test "example" {
    const allocator = std.testing.allocator;
    const input =
        \\L68
        \\L30
        \\R48
        \\L5
        \\R60
        \\L55
        \\L1
        \\L99
        \\R14
        \\L82
    ;
    const instructions = try parseInput(allocator, input);
    defer allocator.free(instructions);

    try std.testing.expectEqual(@as(i32, 3), part1(instructions));
    try std.testing.expectEqual(@as(i32, 6), part2(instructions));
}
