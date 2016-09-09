var webpack = require("webpack");
const DotenvPlugin = require("webpack-dotenv-plugin");
const ExtractTextPlugin = require("extract-text-webpack-plugin");

module.exports = {
  devtool: "source-map",
  entry: {
    "js": "./src/index.jsx",
    "css": "./src/index.scss"
  },
  externals: {
    // "react/addons": true,
    "react/lib/ExecutionEnvironment": true,
    "react/lib/ReactContext": true,
    "assets": true
  },
  module: {
    loaders: [
      {
        test: /\.jsx?$/,
        loader: "babel-loader",
        exclude: /node_modules/,
        query: {
          presets: ["es2015", "react"]
        }
      },
       {
        test: /\.s?css$/,
        loader: ExtractTextPlugin.extract(["style"], ["css", "sass"])
      },
      ,
      {
        test: /\.(woff|woff2|eot|ttf|svg)$/,
        loader: "file?name=[path][name].[ext]"
      }
    ]
  },
  plugins: [
    new ExtractTextPlugin("bundle.css"),
    new DotenvPlugin()
  ],
  output: {
    path: __dirname,
    filename: "./bundle.js"
  },
  resolve: {
    extensions: ["", ".js", ".jsx"]
  }
};
