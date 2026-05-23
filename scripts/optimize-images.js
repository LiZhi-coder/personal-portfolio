import sharp from "sharp";
import { readdir, stat, unlink } from "node:fs/promises";
import { join, extname, basename, dirname } from "node:path";
import { fileURLToPath } from "node:url";
import { existsSync } from "node:fs";

const __dirname = dirname(fileURLToPath(import.meta.url));
const root = join(__dirname, "..");
const publicDir = join(root, "client", "public");
const photosDir = join(publicDir, "images", "photos");
const projectsDir = join(publicDir, "images", "projects");

const PHOTO_QUALITY = 80;
const PNG_QUALITY = 85;
const AVATAR_SIZE = 256;

function formatSize(bytes) {
  if (bytes >= 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  return `${(bytes / 1024).toFixed(1)} KB`;
}

async function convertOne(inputPath, outputPath, options = {}) {
  const { resize, quality = PHOTO_QUALITY } = options;

  const originalSize = (await stat(inputPath)).size;

  let pipeline = sharp(inputPath);
  if (resize) {
    pipeline = pipeline.resize(resize, resize, { fit: "inside", withoutEnlargement: true });
  }
  await pipeline.webp({ quality, effort: 4 }).toFile(outputPath);

  const newSize = (await stat(outputPath)).size;
  const ratio = ((1 - newSize / originalSize) * 100).toFixed(1);

  return { input: basename(inputPath), output: basename(outputPath), originalSize, newSize, ratio };
}

async function processDirectory(dir, quality) {
  let entries;
  try {
    entries = await readdir(dir, { withFileTypes: true });
  } catch {
    return [];
  }

  const results = [];
  const sourceFiles = entries.filter((e) => {
    if (!e.isFile()) return false;
    const ext = extname(e.name).toLowerCase();
    return ext === ".jpg" || ext === ".jpeg" || ext === ".png";
  });

  for (const file of sourceFiles) {
    const ext = extname(file.name);
    const base = basename(file.name, ext);
    const inputPath = join(dir, file.name);
    const outputPath = join(dir, `${base}.webp`);

    // Skip if webp already exists
    if (existsSync(outputPath)) {
      console.log(`  ${file.name}  →  ${base}.webp (skip, already exists)`);
      continue;
    }

    const isPng = ext.toLowerCase() === ".png";
    const result = await convertOne(inputPath, outputPath, {
      quality: isPng ? PNG_QUALITY : quality,
    });

    // If webp is larger, delete it and keep original
    if (result.newSize > result.originalSize) {
      await unlink(outputPath);
      console.log(`  ${result.input}  →  webp larger, keeping original`);
      continue;
    }

    // Success: delete original
    await unlink(inputPath);

    results.push(result);
    console.log(`  ${result.input}  →  ${result.output}`);
    console.log(`    ${formatSize(result.originalSize)} → ${formatSize(result.newSize)} (${result.ratio}%)`);
  }

  return results;
}

async function main() {
  console.log("Optimizing images...\n");

  const results = [];
  let totalOriginal = 0;
  let totalNew = 0;

  // 1. Photos directory
  const photoResults = await processDirectory(photosDir, PHOTO_QUALITY);
  results.push(...photoResults);

  // 2. Project images
  try {
    const projectDirs = await readdir(projectsDir, { withFileTypes: true });
    for (const dir of projectDirs) {
      if (!dir.isDirectory()) continue;
      const subResults = await processDirectory(join(projectsDir, dir.name), PNG_QUALITY);
      results.push(...subResults);
    }
  } catch {
    // projects dir may not exist
  }

  // 3. Avatar (one-time, already converted)
  const avatarPath = join(publicDir, "avatar1.png");
  const avatarWebpPath = join(publicDir, "avatar1.webp");
  if (existsSync(avatarPath) && !existsSync(avatarWebpPath)) {
    console.log("\nProcessing avatar...");
    const avatarResult = await convertOne(avatarPath, avatarWebpPath, {
      resize: AVATAR_SIZE,
      quality: PNG_QUALITY,
    });
    results.push(avatarResult);
    await unlink(avatarPath);
    console.log(
      `  ${avatarResult.input}  →  ${avatarResult.output} (resized to ${AVATAR_SIZE}px)`,
    );
    console.log(`    ${formatSize(avatarResult.originalSize)} → ${formatSize(avatarResult.newSize)} (${avatarResult.ratio}%)`);
  }

  // Summary
  for (const r of results) {
    totalOriginal += r.originalSize;
    totalNew += r.newSize;
  }

  if (results.length > 0) {
    const totalRatio = ((1 - totalNew / totalOriginal) * 100).toFixed(1);
    console.log(`\n${"=".repeat(50)}`);
    console.log(`Total: ${formatSize(totalOriginal)} → ${formatSize(totalNew)} (${totalRatio}%)`);
    console.log(`Optimized ${results.length} images`);
  } else {
    console.log("\nAll images up to date.");
  }
}

main().catch((err) => {
  console.error("Image optimization failed:", err);
  process.exit(1);
});
