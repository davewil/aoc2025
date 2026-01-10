const std = @import("std");
const utils = @import("utils");

const Input = struct {
    operators: [][]const u8,
    expressions: [][]i32,
    rotated: utils.Grid(u8),
    allocator: std.mem.Allocator,

    pub fn deinit(self: *Input) void {
        for (self.operators) |op| self.allocator.free(op);
        self.allocator.free(self.operators);
        for (self.expressions) |expr| self.allocator.free(expr);
        self.allocator.free(self.expressions);
        var mutable_rotated = self.rotated;
        mutable_rotated.deinit();
    }
};

fn parseInput(allocator: std.mem.Allocator, raw: []const u8) !Input {
    var lines_list = std.ArrayListUnmanaged([]const u8){};
    defer lines_list.deinit(allocator);
    var it = std.mem.splitScalar(u8, raw, '\n');
    while (it.next()) |line| {
        if (line.len > 0) {
            try lines_list.append(allocator, line);
        }
    }

    const grid = try utils.Grid(u8).fromLines(allocator, raw);
    defer {
        var mutable_g = grid;
        mutable_g.deinit();
    }
    const rotated = try grid.rotateLeft();

    var data = std.ArrayListUnmanaged([]const u8){};
    defer data.deinit(allocator);
    for (lines_list.items) |line| {
        try data.append(allocator, line);
    }
    std.mem.reverse([]const u8, data.items);

    // operators are in data[0]
    var op_list = std.ArrayListUnmanaged([]const u8){};
    errdefer {
        for (op_list.items) |op| allocator.free(op);
        op_list.deinit(allocator);
    }
    var op_it = std.mem.tokenizeAny(u8, data.items[0], " ");
    while (op_it.next()) |op| {
        try op_list.append(allocator, try allocator.dupe(u8, op));
    }

    var expressions = std.ArrayListUnmanaged([]i32){};
    errdefer {
        for (expressions.items) |expr| allocator.free(expr);
        expressions.deinit(allocator);
    }
    for (data.items[1..]) |line| {
        var expr = std.ArrayListUnmanaged(i32){};
        errdefer expr.deinit(allocator);
        var f_it = std.mem.tokenizeAny(u8, line, " ");
        while (f_it.next()) |field| {
            try expr.append(allocator, try std.fmt.parseInt(i32, field, 10));
        }
        try expressions.append(allocator, try expr.toOwnedSlice(allocator));
    }

    return Input{
        .operators = try op_list.toOwnedSlice(allocator),
        .expressions = try expressions.toOwnedSlice(allocator),
        .rotated = rotated,
        .allocator = allocator,
    };
}

fn applyOp(op: []const u8, values: []const i32) i32 {
    if (std.mem.eql(u8, op, "+")) {
        var sum: i32 = 0;
        for (values) |v| sum += v;
        return sum;
    } else if (std.mem.eql(u8, op, "*")) {
        var prod: i32 = 1;
        for (values) |v| prod *= v;
        return prod;
    }
    return 0;
}

fn part1(allocator: std.mem.Allocator, input: Input) !i32 {
    var total: i32 = 0;
    for (input.operators, 0..) |op, i| {
        var col = std.ArrayListUnmanaged(i32){};
        defer col.deinit(allocator);
        for (input.expressions) |row| {
            if (i < row.len) {
                try col.append(allocator, row[i]);
            }
        }
        total += applyOp(op, col.items);
    }
    return total;
}

fn part2(allocator: std.mem.Allocator, input: Input) !i32 {
    var total: i32 = 0;
    var operands = std.ArrayListUnmanaged(i32){};
    defer operands.deinit(allocator);
    var operator: [1]u8 = undefined;
    var has_operator = false;

    const applyPending = struct {
        fn apply(t: *i32, ops: *std.ArrayListUnmanaged(i32), o: []u8, has_o: *bool) void {
            if (ops.items.len == 0 or !has_o.*) return;
            t.* += applyOp(o, ops.items);
            ops.clearRetainingCapacity();
            has_o.* = false;
        }
    }.apply;

    for (0..input.rotated.rows) |r| {
        const line = input.rotated.row(r);
        const trimmed = std.mem.trim(u8, line, " ");
        if (trimmed.len == 0) {
            applyPending(&total, &operands, &operator, &has_operator);
            continue;
        }

        if (std.mem.indexOfAny(u8, line, "+*")) |idx| {
            const number_part = std.mem.trim(u8, line[0..idx], " ");
            if (number_part.len > 0) {
                if (std.fmt.parseInt(i32, number_part, 10)) |num| {
                    try operands.append(allocator, num);
                } else |_| {}
            }
            operator[0] = line[idx];
            has_operator = true;
            applyPending(&total, &operands, &operator, &has_operator);
            continue;
        }

        if (std.fmt.parseInt(i32, trimmed, 10)) |num| {
            try operands.append(allocator, num);
        } else |_| {}
    }
    applyPending(&total, &operands, &operator, &has_operator);

    return total;
}

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    const raw = try utils.input.getPuzzleInput(allocator, 6);
    defer allocator.free(raw);

    var input = try parseInput(allocator, raw);
    defer input.deinit();

    const p1 = try part1(allocator, input);
    std.debug.print("Part 1: {d}\n", .{p1});

    const p2 = try part2(allocator, input);
    std.debug.print("Part 2: {d}\n", .{p2});
}

test "example" {
    const allocator = std.testing.allocator;
    const raw =
        \\123 328  51 64 
        \\ 45 64  387 23 
        \\  6 98  215 314
        \\*   +   *   +  
    ;
    var input = try parseInput(allocator, raw);
    defer input.deinit();

    try std.testing.expectEqual(@as(i32, 4277556), try part1(allocator, input));
    try std.testing.expectEqual(@as(i32, 3263827), try part2(allocator, input));
}
