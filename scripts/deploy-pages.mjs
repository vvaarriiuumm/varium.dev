import { spawnSync } from "node:child_process";

const projectName = "varium";
const productionBranch = "main";
const distDir = "dist";
const npxCmd = process.platform === "win32" ? "npx.cmd" : "npx";

function run(command, args, options = {}) {
  const result = spawnSync(command, args, {
    stdio: "pipe",
    encoding: "utf8",
    shell: process.platform === "win32",
    ...options,
  });

  return {
    status: result.status ?? 1,
    stdout: result.stdout || "",
    stderr: result.stderr || "",
    output: `${result.stdout || ""}\n${result.stderr || ""}`,
  };
}

// Best-effort create (required only once). If it already exists, continue.
const createResult = run(npxCmd, [
  "wrangler",
  "pages",
  "project",
  "create",
  projectName,
  "--production-branch",
  productionBranch,
]);

if (createResult.status !== 0) {
  const alreadyExists =
    createResult.output.includes("already exists") ||
    createResult.output.includes("already been taken");

  if (!alreadyExists) {
    process.stderr.write(createResult.output);
    process.exit(createResult.status);
  }
}

const deployResult = run(npxCmd, [
  "wrangler",
  "pages",
  "deploy",
  distDir,
  "--project-name",
  projectName,
]);

process.stdout.write(deployResult.stdout);
process.stderr.write(deployResult.stderr);
process.exit(deployResult.status);
