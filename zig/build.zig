const std = @import("std");

pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    // Create the utils module
    const utils_module = b.addModule("utils", .{
        .root_source_file = b.path("src/utils/root.zig"),
    });

    const days = [_][]const u8{ "01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12" };

    // Set up a run step for a specific day: zig build run -Dday=01
    const day_option = b.option([]const u8, "day", "The day to run (e.g., 01)");

    for (days) |day| {
        const name = b.fmt("day{s}", .{day});
        const path = b.fmt("src/day{s}/main.zig", .{day});

        // Check if the file exists before adding it (we'll port them one by one)
        std.fs.cwd().access(path, .{}) catch continue;

        const exe = b.addExecutable(.{
            .name = name,
            .root_module = b.createModule(.{
                .root_source_file = b.path(path),
                .target = target,
                .optimize = optimize,
            }),
        });
        exe.root_module.addImport("utils", utils_module);

        // Day 10 needs Z3
        if (std.mem.eql(u8, day, "10")) {
            exe.linkLibC();
            exe.linkSystemLibrary("z3");
        }

        b.installArtifact(exe);

        const run_cmd = b.addRunArtifact(exe);
        run_cmd.step.dependOn(b.getInstallStep());

        if (day_option) |d| {
            if (std.mem.eql(u8, d, day)) {
                const run_step = b.step("run", b.fmt("Run day {s}", .{day}));
                run_step.dependOn(&run_cmd.step);
            }
        }

        const exe_unit_tests = b.addTest(.{
            .root_module = b.createModule(.{
                .root_source_file = b.path(path),
                .target = target,
                .optimize = optimize,
            }),
        });
        exe_unit_tests.root_module.addImport("utils", utils_module);

        if (std.mem.eql(u8, day, "10")) {
            exe_unit_tests.linkLibC();
            exe_unit_tests.linkSystemLibrary("z3");
        }

        const run_exe_unit_tests = b.addRunArtifact(exe_unit_tests);

        const test_step = b.step(b.fmt("test_{s}", .{day}), b.fmt("Run unit tests for day {s}", .{day}));
        test_step.dependOn(&run_exe_unit_tests.step);
    }
}
