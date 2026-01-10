const std = @import("std");
const utils = @import("utils");
const z3 = @cImport({
    @cInclude("z3.h");
});

const Row = struct {
    lights: []usize,
    buttons: [][]usize,
    joltages: []i32,
    allocator: std.mem.Allocator,

    pub fn deinit(self: *Row) void {
        self.allocator.free(self.lights);
        for (self.buttons) |b| self.allocator.free(b);
        self.allocator.free(self.buttons);
        self.allocator.free(self.joltages);
    }
};

const Input = struct {
    rows: []Row,
    allocator: std.mem.Allocator,

    pub fn deinit(self: *Input) void {
        for (self.rows) |*row| row.deinit();
        self.allocator.free(self.rows);
    }
};

fn parseRow(allocator: std.mem.Allocator, line: []const u8) !Row {
    var lights = std.ArrayListUnmanaged(usize){};
    errdefer lights.deinit(allocator);
    var buttons = std.ArrayListUnmanaged([]usize){};
    errdefer {
        for (buttons.items) |b| allocator.free(b);
        buttons.deinit(allocator);
    }
    var joltages = std.ArrayListUnmanaged(i32){};
    errdefer joltages.deinit(allocator);

    const l_start = std.mem.indexOfScalar(u8, line, '[') orelse return error.InvalidFormat;
    const l_end = std.mem.indexOfScalar(u8, line, ']') orelse return error.InvalidFormat;
    const lights_str = line[l_start + 1 .. l_end];
    for (lights_str, 0..) |c, i| {
        if (c == '#') try lights.append(allocator, i);
    }

    var pos = l_end + 1;
    while (pos < line.len) {
        const next_open = std.mem.indexOfAnyPos(u8, line, pos, "({") orelse break;
        const next_close = if (line[next_open] == '(')
            std.mem.indexOfScalarPos(u8, line, next_open, ')')
        else
            std.mem.indexOfScalarPos(u8, line, next_open, '}');
        
        const end = next_close orelse break;
        const content = line[next_open + 1 .. end];
        
        if (line[next_open] == '(') {
            var b_list = std.ArrayListUnmanaged(usize){};
            errdefer b_list.deinit(allocator);
            var it = std.mem.tokenizeAny(u8, content, ", ");
            while (it.next()) |token| {
                const t = std.mem.trim(u8, token, " \t\r\n");
                if (t.len == 0) continue;
                try b_list.append(allocator, try std.fmt.parseInt(usize, t, 10));
            }
            try buttons.append(allocator, try b_list.toOwnedSlice(allocator));
        } else {
            var it = std.mem.tokenizeAny(u8, content, ", ");
            while (it.next()) |token| {
                const t = std.mem.trim(u8, token, " \t\r\n");
                if (t.len == 0) continue;
                try joltages.append(allocator, try std.fmt.parseInt(i32, t, 10));
            }
        }
        pos = end + 1;
    }

    return Row{
        .lights = try lights.toOwnedSlice(allocator),
        .buttons = try buttons.toOwnedSlice(allocator),
        .joltages = try joltages.toOwnedSlice(allocator),
        .allocator = allocator,
    };
}

fn parseInput(allocator: std.mem.Allocator, raw: []const u8) !Input {
    var rows = std.ArrayListUnmanaged(Row){};
    errdefer {
        for (rows.items) |*r| r.deinit();
        rows.deinit(allocator);
    }

    var it = std.mem.tokenizeAny(u8, raw, "\n\r");
    while (it.next()) |line| {
        if (std.mem.trim(u8, line, " ").len == 0) continue;
        try rows.append(allocator, try parseRow(allocator, line));
    }

    return Input{
        .rows = try rows.toOwnedSlice(allocator),
        .allocator = allocator,
    };
}

fn rowLength(row: Row) usize {
    var max: usize = 0;
    for (row.lights) |idx| {
        if (idx + 1 > max) max = idx + 1;
    }
    for (row.buttons) |btn| {
        for (btn) |idx| {
            if (idx + 1 > max) max = idx + 1;
        }
    }
    return max;
}

fn part1(allocator: std.mem.Allocator, input: Input) !i32 {
    var total: i32 = 0;
    for (input.rows) |row| {
        total += try findMinPresses(allocator, row);
    }
    return total;
}

fn findMinPresses(allocator: std.mem.Allocator, row: Row) !i32 {
    const n = rowLength(row);
    if (n == 0) return 0;
    
    var target: u64 = 0;
    for (row.lights) |idx| {
        if (idx < 64) target |= (@as(u64, 1) << @as(u6, @intCast(idx)));
    }
    if (target == 0) return 0;

    var button_masks = try allocator.alloc(u64, row.buttons.len);
    defer allocator.free(button_masks);
    for (row.buttons, 0..) |btn, i| {
        var m: u64 = 0;
        for (btn) |idx| {
            if (idx < 64) m |= (@as(u64, 1) << @as(u6, @intCast(idx)));
        }
        button_masks[i] = m;
    }

    const max_states = @as(usize, 1) << @as(u6, @intCast(n));
    var dist = try allocator.alloc(i32, max_states);
    defer allocator.free(dist);
    @memset(dist, -1);

    var q = std.ArrayListUnmanaged(u64){};
    defer q.deinit(allocator);

    try q.append(allocator, 0);
    dist[0] = 0;

    var head: usize = 0;
    while (head < q.items.len) {
        const s = q.items[head];
        head += 1;

        if (s == target) return dist[@intCast(s)];

        for (button_masks) |bm| {
            const ns = s ^ bm;
            if (ns < max_states and dist[@intCast(ns)] == -1) {
                dist[@intCast(ns)] = dist[@intCast(s)] + 1;
                try q.append(allocator, ns);
            }
        }
    }

    return -1;
}

fn findMinPressesWithZ3(row: Row) !i32 {
    if (row.joltages.len == 0) return 0;
    if (row.buttons.len == 0) {
        for (row.joltages) |v| if (v != 0) return -1;
        return 0;
    }

    const cfg = z3.Z3_mk_config();
    defer z3.Z3_del_config(cfg);
    const ctx = z3.Z3_mk_context(cfg);
    defer z3.Z3_del_context(ctx);
    const slv = z3.Z3_mk_solver(ctx);
    z3.Z3_solver_inc_ref(ctx, slv);
    defer z3.Z3_solver_dec_ref(ctx, slv);

    const int_sort = z3.Z3_mk_int_sort(ctx);
    const zero = z3.Z3_mk_int64(ctx, 0, int_sort);

    const n_buttons = row.buttons.len;
    const n_lights = row.joltages.len;

    var button_vars = try row.allocator.alloc(z3.Z3_ast, n_buttons);
    defer row.allocator.free(button_vars);

    for (0..n_buttons) |i| {
        var buf: [32]u8 = undefined;
        const name = try std.fmt.bufPrintZ(&buf, "b{d}", .{i});
        const sym = z3.Z3_mk_string_symbol(ctx, name);
        button_vars[i] = z3.Z3_mk_const(ctx, sym, int_sort);
        z3.Z3_solver_assert(ctx, slv, z3.Z3_mk_ge(ctx, button_vars[i], zero));
    }

    for (0..n_lights) |l_idx| {
        var terms = std.ArrayListUnmanaged(z3.Z3_ast){};
        defer terms.deinit(row.allocator);
        for (row.buttons, 0..) |btn, b_idx| {
            var affects = false;
            for (btn) |idx| if (idx == l_idx) { affects = true; break; };
            if (affects) try terms.append(row.allocator, button_vars[b_idx]);
        }

        if (terms.items.len == 0) {
            if (row.joltages[l_idx] != 0) return -1;
            continue;
        }

        const sum = if (terms.items.len == 1) terms.items[0] else z3.Z3_mk_add(ctx, @intCast(terms.items.len), terms.items.ptr);
        z3.Z3_solver_assert(ctx, slv, z3.Z3_mk_eq(ctx, sum, z3.Z3_mk_int64(ctx, row.joltages[l_idx], int_sort)));
    }

    const total_presses = if (button_vars.len == 1) button_vars[0] else z3.Z3_mk_add(ctx, @intCast(button_vars.len), button_vars.ptr);

    var max_sum: i32 = 0;
    for (row.joltages) |j| max_sum += j;

    var max_k: usize = 0;
    for (row.buttons) |btn| if (btn.len > max_k) { max_k = btn.len; };

    var lo: i32 = 0;
    if (max_k > 0) lo = @divTrunc(max_sum + @as(i32, @intCast(max_k)) - 1, @as(i32, @intCast(max_k)));
    
    var hi = max_sum;
    var best: i32 = -1;

    while (lo <= hi) {
        const mid = @divTrunc(lo + hi, 2);
        z3.Z3_solver_push(ctx, slv);
        z3.Z3_solver_assert(ctx, slv, z3.Z3_mk_le(ctx, total_presses, z3.Z3_mk_int64(ctx, mid, int_sort)));
        
        const res = z3.Z3_solver_check(ctx, slv);
        if (res == z3.Z3_L_TRUE) {
            const model = z3.Z3_solver_get_model(ctx, slv);
            z3.Z3_model_inc_ref(ctx, model);
            var val_ast: z3.Z3_ast = undefined;
            _ = z3.Z3_model_eval(ctx, model, total_presses, true, &val_ast);
            
            var val: i64 = 0;
            _ = z3.Z3_get_numeral_int64(ctx, val_ast, &val);
            z3.Z3_model_dec_ref(ctx, model);

            best = @intCast(val);
            hi = best - 1;
        } else {
            lo = mid + 1;
        }
        z3.Z3_solver_pop(ctx, slv, 1);
    }

    return best;
}

fn part2(input: Input) !i32 {
    var total: i32 = 0;
    for (input.rows) |row| {
        const v = try findMinPressesWithZ3(row);
        if (v == -1) return -1;
        total += v;
    }
    return total;
}

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    const raw = try utils.input.getPuzzleInput(allocator, 10);
    defer allocator.free(raw);

    var input = try parseInput(allocator, raw);
    defer input.deinit();

    const p1 = try part1(allocator, input);
    std.debug.print("Part 1: {d}\n", .{p1});

    const p2 = try part2(input);
    std.debug.print("Part 2: {d}\n", .{p2});
}

test "example" {
    const allocator = std.testing.allocator;
            const raw = "[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}\n[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}\n[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}";
            var input = try parseInput(allocator, raw);
            defer input.deinit();
        
            try std.testing.expectEqual(@as(i32, 7), try part1(allocator, input));
            try std.testing.expectEqual(@as(i32, 33), try part2(input));
        }
        
    