import { exists, rm } from "fs/promises";
import tailwind from "bun-plugin-tailwind";

const outdir = "dist";

if (await exists(outdir)) {
  console.log(`🗑️ Cleaning previous build at ${outdir}`);
  await rm(outdir, { recursive: true, force: true });
}

const result = await Bun.build({
  entrypoints: ["src/index.html"],
  outdir,
  plugins: [tailwind],
  minify: true,
  target: "browser",
  sourcemap: "linked",
  publicPath: "dist/",
});

console.log(result);
