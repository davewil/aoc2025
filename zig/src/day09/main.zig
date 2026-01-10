const std = @import("std");
const utils = @import("utils");

const Point2D = struct {
    x: i32,
    y: i32,

    pub fn eql(self: Point2D, other: Point2D) bool {
        return self.x == other.x and self.y == other.y;
    }
};

const Edge = struct {
    from: Point2D,
    to: Point2D,
};

const Rectangle = struct {
    min_x: i32,
    min_y: i32,
    max_x: i32,
    max_y: i32,
};

fn parseInput(allocator: std.mem.Allocator, raw: []const u8) ![]Point2D {
    var coords = std.ArrayListUnmanaged(Point2D){};
    errdefer coords.deinit(allocator);

    var it = std.mem.tokenizeAny(u8, raw, "\n\r");
    while (it.next()) |line| {
        if (std.mem.indexOfScalar(u8, line, ',')) |idx| {
            const before = std.mem.trim(u8, line[0..idx], " ");
            const after = std.mem.trim(u8, line[idx + 1 ..], " ");
            try coords.append(allocator, .{
                .x = try std.fmt.parseInt(i32, before, 10),
                .y = try std.fmt.parseInt(i32, after, 10),
            });
        }
    }
    return coords.toOwnedSlice(allocator);
}

fn part1(coords: []const Point2D) i32 {
    var max_area: i32 = 0;
    for (0..coords.len) |from| {
        for (from + 1..coords.len) |to| {
            const p1 = coords[from];
            const p2 = coords[to];
            const area = (utils.abs(p1.x - p2.x) + 1) * (utils.abs(p1.y - p2.y) + 1);
            if (area > max_area) max_area = area;
        }
    }
    return max_area;
}

fn buildEdges(allocator: std.mem.Allocator, red_tiles: []const Point2D) ![]Edge {
    const edges = try allocator.alloc(Edge, red_tiles.len);
    for (red_tiles, 0..) |from_point, i| {
        const to_point = red_tiles[(i + 1) % red_tiles.len];
        edges[i] = .{ .from = from_point, .to = to_point };
    }
    return edges;
}

fn sort(a: i32, b: i32) struct { i32, i32 } {
    if (a < b) return .{ a, b };
    return .{ b, a };
}

fn strictOverlap(a_min: i32, a_max: i32, b_min: i32, b_max: i32) bool {
    return a_max > b_min and b_max > a_min;
}

fn intersects(edges: []const Edge, from: Point2D, to: Point2D) bool {
    const s_x = sort(from.x, to.x);
    const min_x = s_x[0];
    const max_x = s_x[1];
    const s_y = sort(from.y, to.y);
    const min_y = s_y[0];
    const max_y = s_y[1];

    for (edges) |edge| {
        if (edge.from.x == edge.to.x) {
            const x = edge.from.x;
            if (x <= min_x or x >= max_x) continue;
            const e_y = sort(edge.from.y, edge.to.y);
            if (strictOverlap(min_y, max_y, e_y[0], e_y[1])) return true;
        } else {
            const y = edge.from.y;
            if (y <= min_y or y >= max_y) continue;
            const e_x = sort(edge.from.x, edge.to.x);
            if (strictOverlap(min_x, max_x, e_x[0], e_x[1])) return true;
        }
    }
    return false;
}

fn hasInteriorVertex(vertices: []const Point2D, min_x: i32, max_x: i32, min_y: i32, max_y: i32, skip_a: Point2D, skip_b: Point2D) bool {
    for (vertices) |p| {
        if (p.eql(skip_a) or p.eql(skip_b)) continue;
        if (p.x > min_x and p.x < max_x and p.y > min_y and p.y < max_y) return true;
    }
    return false;
}

fn pointInPolygonCenter(x2: i32, y2: i32, polygon: []const Point2D) bool {
    var inside = false;
    const n = polygon.len;
    for (0..n) |i| {
        const j = (i + n - 1) % n;
        const yi2 = polygon[i].y * 2;
        const yj2 = polygon[j].y * 2;
        if ((yi2 > y2) == (yj2 > y2)) continue;

        const xi2 = polygon[i].x * 2;
        const xj2 = polygon[j].x * 2;
        const dy = yj2 - yi2;
        if (dy == 0) continue;

        const dx = xj2 - xi2;
        const y_delta = y2 - yi2;
        const lhs = @as(i64, x2) * dy;
        const rhs = @as(i64, xi2) * dy + @as(i64, dx) * y_delta;

        if ((dy > 0 and lhs < rhs) or (dy < 0 and lhs > rhs)) {
            inside = !inside;
        }
    }
    return inside;
}

fn part2(allocator: std.mem.Allocator, red_tiles: []const Point2D) !i32 {
    if (red_tiles.len == 0) return 0;
    const edges = try buildEdges(allocator, red_tiles);
    defer allocator.free(edges);

    var inside_cache = std.AutoHashMap(i64, bool).init(allocator);
    defer inside_cache.deinit();

    var max_area: i32 = 0;
    for (0..red_tiles.len) |from| {
        for (from + 1..red_tiles.len) |to| {
            const p1 = red_tiles[from];
            const p2 = red_tiles[to];

            const s_x = sort(p1.x, p2.x);
            const min_x = s_x[0];
            const max_x = s_x[1];
            const s_y = sort(p1.y, p2.y);
            const min_y = s_y[0];
            const max_y = s_y[1];

            if (min_x == max_x or min_y == max_y) continue;

            if (intersects(edges, p1, p2)) continue;

            const center_sum_x = min_x + max_x;
            const center_sum_y = min_y + max_y;
            const cache_key = (@as(i64, center_sum_x) << 32) | @as(i64, @as(u32, @bitCast(center_sum_y)));
            
            var inside: bool = undefined;
            if (inside_cache.get(cache_key)) |v| {
                inside = v;
            } else {
                inside = pointInPolygonCenter(center_sum_x, center_sum_y, red_tiles);
                try inside_cache.put(cache_key, inside);
            }

            if (!inside) continue;

            if (hasInteriorVertex(red_tiles, min_x, max_x, min_y, max_y, p1, p2)) continue;

            const area = (max_x - min_x + 1) * (max_y - min_y + 1);
            if (area > max_area) {
                max_area = area;
            }
        }
    }
    return max_area;
}

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    const raw = try utils.input.getPuzzleInput(allocator, 9);
    defer allocator.free(raw);

    const coords = try parseInput(allocator, raw);
    defer allocator.free(coords);

    const p1 = part1(coords);
    std.debug.print("Part 1: {d}\n", .{p1});

    const p2 = try part2(allocator, coords);
    std.debug.print("Part 2: {d}\n", .{p2});
}

test "example" {
    const allocator = std.testing.allocator;
    const raw =
        \\7,1
        \\11,1
        \\11,7
        \\9,7
        \\9,5
        \\2,5
        \\2,3
        \\7,3
    ;
    const coords = try parseInput(allocator, raw);
    defer allocator.free(coords);

    try std.testing.expectEqual(@as(i32, 50), part1(coords));
    try std.testing.expectEqual(@as(i32, 24), try part2(allocator, coords));
}
