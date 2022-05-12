const BundleAnalyzerPlugin = require("webpack-bundle-analyzer").BundleAnalyzerPlugin;
const { NormalModuleReplacementPlugin } = require("webpack");
const CopyPlugin = require("copy-webpack-plugin");

module.exports = {
  entry: "./src/index.ts",
  mode: "development",
  devtool: false,
  module: {
    // Use `ts-loader` on any file that ends in '.ts'
    rules: [
      {
        test: /\.ts$/,
        use: "ts-loader",
        exclude: /node_modules/,
      },
    ],
  },
  resolve: {
    extensions: [".ts", ".js"],
  },
  output: {
    filename: "editor.js",
    path: `${process.cwd()}/dist`,
  },
  optimization: {
    // minimize: true,
  },
  plugins: [
    new CopyPlugin({
      patterns: [
        { from: "./assets", to: "" },
      ],
    }),
    // new BundleAnalyzerPlugin({ openAnalyzer: false }),
    // new NormalModuleReplacementPlugin(
    //   /highlight\.js\/lib\/core/,
    //   __dirname + '/src/highlight.ts'
    // ),
  ],
  externals: {
    'highlight.js/lib/core': 'hljs',
    // 'highlight.js': 'hljs',
    // 'highlight.js/lib/languages/php': 'hljs.getLanguage(\'php\')',
    // 'highlight.js/lib/languages/typescript': 'hljs.getLanguage(\'typescript\')',
  },
};
