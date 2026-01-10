const std = @import("std");
const utils = @import("utils");

const Graph = struct {
    adj: [][]usize,
    nodes: std.StringHashMap(usize),
    allocator: std.mem.Allocator,

    pub fn deinit(self: *Graph) void {
        for (self.adj) |neighbors| self.allocator.free(neighbors);
        self.allocator.free(self.adj);
        var it = self.nodes.keyIterator();
        while (it.next()) |name| self.allocator.free(name.*);
        self.nodes.deinit();
    }
};

fn parseInput(allocator: std.mem.Allocator, raw: []const u8) !Graph {
    var nodes = std.StringHashMap(usize).init(allocator);
    errdefer {
        var it = nodes.keyIterator();
        while (it.next()) |name| allocator.free(name.*);
        nodes.deinit();
    }

    var next_id: usize = 0;
    const getID = struct {
        fn get(n: *std.StringHashMap(usize), nid: *usize, name: []const u8, alloc: std.mem.Allocator) !usize {
            if (n.get(name)) |id| return id;
            const id = nid.*;
            nid.* += 1;
            try n.put(try alloc.dupe(u8, name), id);
            return id;
        }
    }.get;

    const Edge = struct { from: usize, to: usize };
    var edges = std.ArrayListUnmanaged(Edge){};
    defer edges.deinit(allocator);

    var it = std.mem.tokenizeAny(u8, raw, "\n\r");
    while (it.next()) |line| {
        var parts = std.mem.splitSequence(u8, line, ": ");
        const from_name = parts.next() orelse continue;
        const children_str = parts.next() orelse continue;

        const from_id = try getID(&nodes, &next_id, std.mem.trim(u8, from_name, " "), allocator);
        var c_it = std.mem.tokenizeAny(u8, children_str, " ");
        while (c_it.next()) |child| {
            const to_id = try getID(&nodes, &next_id, std.mem.trim(u8, child, " "), allocator);
            try edges.append(allocator, .{ .from = from_id, .to = to_id });
        }
    }

    var adj = try allocator.alloc([]usize, next_id);
    for (0..next_id) |i| adj[i] = &[_]usize{};
    
    var adj_builders = try allocator.alloc(std.ArrayListUnmanaged(usize), next_id);
    defer allocator.free(adj_builders);
    for (0..next_id) |i| adj_builders[i] = std.ArrayListUnmanaged(usize){};

    for (edges.items) |e| {
        try adj_builders[e.from].append(allocator, e.to);
    }

    for (0..next_id) |i| {
        adj[i] = try adj_builders[i].toOwnedSlice(allocator);
    }

    return Graph{
        .adj = adj,
        .nodes = nodes,
        .allocator = allocator,
    };
}

fn countPaths(adj: [][]usize, current: usize, target: usize, memo: []i64) i64 {
    if (current == target) return 1;
    if (memo[current] != -1) return memo[current];

    var total: i64 = 0;
    for (adj[current]) |neighbor| {
        total += countPaths(adj, neighbor, target, memo);
    }
    memo[current] = total;
    return total;
}

fn part1(graph: Graph) !i64 {
    const start = graph.nodes.get("you") orelse return 0;
    const end = graph.nodes.get("out") orelse return 0;

    const memo = try graph.allocator.alloc(i64, graph.adj.len);
    defer graph.allocator.free(memo);
    @memset(memo, -1);

    return countPaths(graph.adj, start, end, memo);
}

fn part2(graph: Graph) !i64 {
    const svr = graph.nodes.get("svr") orelse return 0;
    const dac = graph.nodes.get("dac") orelse return 0;
    const fft = graph.nodes.get("fft") orelse return 0;
    const out = graph.nodes.get("out") orelse return 0;

    const memo = try graph.allocator.alloc(i64, graph.adj.len);
    defer graph.allocator.free(memo);

    const runCount = struct {
        fn run(adj: [][]usize, from: usize, to: usize, m: []i64) i64 {
            @memset(m, -1);
            return countPaths(adj, from, to, m);
        }
    }.run;

    const dacToFft = runCount(graph.adj, dac, fft, memo);

    if (dacToFft > 0) {
        const p1 = runCount(graph.adj, svr, dac, memo);
        const p2 = dacToFft;
        const p3 = runCount(graph.adj, fft, out, memo);
        return p1 * p2 * p3;
    } else {
        const fftToDac = runCount(graph.adj, fft, dac, memo);
        if (fftToDac > 0) {
            const p1 = runCount(graph.adj, svr, fft, memo);
            const p2 = fftToDac;
            const p3 = runCount(graph.adj, dac, out, memo);
            return p1 * p2 * p3;
        }
    }

    return 0;
}

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    const raw = try utils.input.getPuzzleInput(allocator, 11);
    defer allocator.free(raw);

    var graph = try parseInput(allocator, raw);
    defer graph.deinit();

    const p1 = try part1(graph);
    std.debug.print("Part 1: {d}\n", .{p1});

    const p2 = try part2(graph);
    std.debug.print("Part 2: {d}\n", .{p2});
}

test "example1" {
    const allocator = std.testing.allocator;
    const raw =
        \\aaa: you hhh
        \\you: bbb ccc
        \\bbb: ddd eee
        \\ccc: ddd eee fff
        \\ddd: ggg
        \\eee: out
        \\fff: out
        \\ggg: out
        \\hhh: ccc fff iii
        \\iii: out
    ;
    var graph = try parseInput(allocator, raw);
    defer graph.deinit();

    try std.testing.expectEqual(@as(i64, 5), try part1(graph));
}

test "example2" {
    const allocator = std.testing.allocator;
    const raw =
        \\svr: aaa bbb
        \\aaa: fft
        \\fft: ccc
        \\bbb: tty
        \\tty: ccc
        \\ccc: ddd eee
        \\ddd: hub
        \\hub: fff
        \\eee: dac
        \\dac: fff
        \\fff: ggg hhh
        \\ggg: out
        \\hhh: out
    ;
    var graph = try parseInput(allocator, raw);
    defer graph.deinit();

    try std.testing.expectEqual(@as(i64, 2), try part2(graph));
}
