const std = @import("std");
const utils = @import("utils");

const Point3D = struct {
    x: i32,
    y: i32,
    z: i32,
};

const Edge = struct {
    a: usize,
    b: usize,
    w: i64,
};

const JunctionSet = struct {
    points: []Point3D,
    edges: []Edge,
    allocator: std.mem.Allocator,

    pub fn deinit(self: *JunctionSet) void {
        self.allocator.free(self.points);
        self.allocator.free(self.edges);
    }
};

const UnionFind = struct {
    parent: []usize,
    rank: []usize,
    allocator: std.mem.Allocator,

    pub fn init(allocator: std.mem.Allocator, size: usize) !UnionFind {
        const parent = try allocator.alloc(usize, size);
        const rank = try allocator.alloc(usize, size);
        for (0..size) |i| {
            parent[i] = i;
            rank[i] = 0;
        }
        return UnionFind{
            .parent = parent,
            .rank = rank,
            .allocator = allocator,
        };
    }

    pub fn deinit(self: *UnionFind) void {
        self.allocator.free(self.parent);
        self.allocator.free(self.rank);
    }

    pub fn find(self: *UnionFind, x: usize) usize {
        if (self.parent[x] != x) {
            self.parent[x] = self.find(self.parent[x]);
        }
        return self.parent[x];
    }

    pub fn unionSet(self: *UnionFind, a: usize, b: usize) bool {
        const rootA = self.find(a);
        const rootB = self.find(b);
        if (rootA == rootB) return false;

        if (self.rank[rootA] < self.rank[rootB]) {
            self.parent[rootA] = rootB;
        } else if (self.rank[rootA] > self.rank[rootB]) {
            self.parent[rootB] = rootA;
        } else {
            self.parent[rootB] = rootA;
            self.rank[rootA] += 1;
        }
        return true;
    }
};

fn squaredDistance(a: Point3D, b: Point3D) i64 {
    const dx = @as(i64, a.x) - b.x;
    const dy = @as(i64, a.y) - b.y;
    const dz = @as(i64, a.z) - b.z;
    return dx * dx + dy * dy + dz * dz;
}

fn buildEdges(allocator: std.mem.Allocator, points: []const Point3D) ![]Edge {
    const n = points.len;
    if (n < 2) return &[_]Edge{};
    var edges = std.ArrayListUnmanaged(Edge){};
    errdefer edges.deinit(allocator);

    for (0..n) |i| {
        for (i + 1..n) |j| {
            try edges.append(allocator, .{
                .a = i,
                .b = j,
                .w = squaredDistance(points[i], points[j]),
            });
        }
    }

    std.mem.sort(Edge, edges.items, {}, struct {
        fn lessThan(_: void, a: Edge, b: Edge) bool {
            return a.w < b.w;
        }
    }.lessThan);

    return edges.toOwnedSlice(allocator);
}

fn parseInput(allocator: std.mem.Allocator, raw: []const u8) !JunctionSet {
    var points = std.ArrayListUnmanaged(Point3D){};
    errdefer points.deinit(allocator);

    var it = std.mem.tokenizeAny(u8, raw, "\n\r");
    while (it.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " ");
        if (trimmed.len == 0) continue;

        var parts = std.mem.splitScalar(u8, trimmed, ',');
        const x_str = parts.next() orelse continue;
        const y_str = parts.next() orelse continue;
        const z_str = parts.next() orelse continue;

        try points.append(allocator, .{
            .x = try std.fmt.parseInt(i32, std.mem.trim(u8, x_str, " "), 10),
            .y = try std.fmt.parseInt(i32, std.mem.trim(u8, y_str, " "), 10),
            .z = try std.fmt.parseInt(i32, std.mem.trim(u8, z_str, " "), 10),
        });
    }

    const pts = try points.toOwnedSlice(allocator);
    const eds = try buildEdges(allocator, pts);

    return JunctionSet{
        .points = pts,
        .edges = eds,
        .allocator = allocator,
    };
}

fn part1(allocator: std.mem.Allocator, data: JunctionSet, connect_count: usize) !usize {
    if (data.points.len == 0) return 0;
    const n = data.points.len;
    var uf = try UnionFind.init(allocator, n);
    defer uf.deinit();

    const actual_connect = @min(connect_count, data.edges.len);
    for (0..actual_connect) |i| {
        const e = data.edges[i];
        _ = uf.unionSet(e.a, e.b);
    }

    var counts = std.AutoHashMap(usize, usize).init(allocator);
    defer counts.deinit();
    for (0..n) |i| {
        const root = uf.find(i);
        const res = try counts.getOrPut(root);
        if (!res.found_existing) res.value_ptr.* = 0;
        res.value_ptr.* += 1;
    }

    var sizes = std.ArrayListUnmanaged(usize){};
    defer sizes.deinit(allocator);
    var vit = counts.valueIterator();
    while (vit.next()) |v| {
        try sizes.append(allocator, v.*);
    }

    std.mem.sort(usize, sizes.items, {}, std.sort.desc(usize));
    const limit = @min(sizes.items.len, 3);
    var product: usize = 1;
    for (0..limit) |i| {
        product *= sizes.items[i];
    }
    return if (sizes.items.len == 0) 0 else product;
}

fn part2(allocator: std.mem.Allocator, data: JunctionSet) !i64 {
    if (data.points.len == 0) return 0;
    const n = data.points.len;
    var uf = try UnionFind.init(allocator, n);
    defer uf.deinit();

    var components = n;
    for (data.edges) |e| {
        if (uf.unionSet(e.a, e.b)) {
            components -= 1;
            if (components == 1) {
                return @as(i64, data.points[e.a].x) * data.points[e.b].x;
            }
        }
    }
    return 0;
}

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    const raw = try utils.input.getPuzzleInput(allocator, 8);
    defer allocator.free(raw);

    var junctions = try parseInput(allocator, raw);
    defer junctions.deinit();

    const p1 = try part1(allocator, junctions, 1000);
    std.debug.print("Part 1: {d}\n", .{p1});

    const p2 = try part2(allocator, junctions);
    std.debug.print("Part 2: {d}\n", .{p2});
}

test "example" {
    const allocator = std.testing.allocator;
    const raw =
        \\162,817,812
        \\57,618,57
        \\906,360,560
        \\592,479,940
        \\352,342,300
        \\466,668,158
        \\542,29,236
        \\431,825,988
        \\739,650,466
        \\52,470,668
        \\216,146,977
        \\819,987,18
        \\117,168,530
        \\805,96,715
        \\346,949,466
        \\970,615,88
        \\941,993,340
        \\862,61,35
        \\984,92,344
        \\425,690,689
    ;
    var junctions = try parseInput(allocator, raw);
    defer junctions.deinit();

    try std.testing.expectEqual(@as(usize, 40), try part1(allocator, junctions, 10));
    try std.testing.expectEqual(@as(i64, 25272), try part2(allocator, junctions));
}
