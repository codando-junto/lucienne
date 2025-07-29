import "dotenv/config"
import * as esbuild from 'esbuild'
import { sassPlugin } from 'esbuild-sass-plugin'
import chokidar from 'chokidar'
import fs from 'fs'

const ENVIRONMENT = process.env.APP_ENV
const ASSETS_PATH = "assets"
const PUBLIC_PATH = "public"
const COMPILED_ASSETS_PATH = `${PUBLIC_PATH}/assets`
const ASSETS_BUILD_FILE = "build.json"
const SUPPORTED_MEDIA_FORMATS = ['.jpg', '.jpeg', '.png', '.ico']
const ENTRYPOINT_PATHS = ["javascript/*.js", "scss/*.scss", "images/**/*"]
const ENTRYNAME_FORMAT = "[dir]/[name]-[hash]"
const SOURCEMAP_TYPE = "linked"

const imagesCopyLoader = SUPPORTED_MEDIA_FORMATS.reduce(
  (object, media) => ({ ...object, [media]: 'copy' }),
  {}
)

const entryPoints = ENTRYPOINT_PATHS.map(path => `${ASSETS_PATH}/${path}`)

let esBuildOptions = {
  entryPoints: entryPoints,
  entryNames: ENTRYNAME_FORMAT,
  outdir: COMPILED_ASSETS_PATH,
  bundle: true,
  minify: true,
  metafile: true,
  loader: imagesCopyLoader,
  plugins: [
    sassPlugin({
      filter: /\.s[ac]ss|css$/,
      loadPaths: [`${ASSETS_PATH}/scss`]
    })
  ],
}

if (ENVIRONMENT == "development") {
  await buildDevelopmentAssets();
} else {
  await buildProductionAssets();
}

async function buildDevelopmentAssets() {
  esBuildOptions.sourcemap = SOURCEMAP_TYPE;
  const context = await esbuild.context(esBuildOptions);

  console.log(`Watching ./${ASSETS_PATH} directory`)

  chokidar.watch(`./${ASSETS_PATH}`).on('all', async (eventType, filename) => {
    buildAndMapAssets(context);
  });

  buildAndMapAssets(context);
  console.log("Assets built")
}

async function buildProductionAssets() {
  const context = await esbuild.context(esBuildOptions);
  buildAndMapAssets(context);
  console.log("Assets built for: Production")
}

async function buildAndMapAssets(context) {
  const result = await context.rebuild();
  const { outputs } = result.metafile;
  const builtAssetsMapping = Object.keys(outputs)
    .filter(key => !key.endsWith(".map"))
    .reduce(
      (obj, key) => {
        const entryPoint = outputs[key].entryPoint
        const builtAsset = key
        return { ...obj, [entryPoint]: builtAsset }
      },
      {}
    )

  fs.writeFile(`${PUBLIC_PATH}/${ASSETS_BUILD_FILE}`, JSON.stringify(builtAssetsMapping), err => {
    if (err) throw err;
  });
}
