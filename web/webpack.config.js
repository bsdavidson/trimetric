var webpack = require("webpack");
const DotenvPlugin = require("webpack-dotenv-plugin");
const ExtractTextPlugin = require("extract-text-webpack-plugin");

module.exports = {
  devtool: "source-map",
  devServer: {
    historyApiFallback: true,
    proxy: {
      "/api": {
        target: "http://api:80"
      }
    }
  },
  entry: ["./src/index.jsx", "./src/index.scss"],
  externals: {
    // "react/addons": true,
    "react/lib/ExecutionEnvironment": true,
    "react/lib/ReactContext": true,
    assets: true
  },
  module: {
    rules: [
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        use: {
          loader: "babel-loader",
          options: {
            presets: ["es2015", "react"]
          }
        }
      },
      {
        test: /\.s?css$/,
        use: ExtractTextPlugin.extract({
          fallback: "style-loader",
          use: ["css-loader", "sass-loader"]
        })
      },
      {
        test: /\.(woff|woff2|eot|ttf|svg)$/,
        use: {
          loader: "file-loader",
          options: {
            name: "[path][name].[ext]"
          }
        }
      }
    ]
  },
  // },
  plugins: [
    new ExtractTextPlugin({filename: "bundle.css", allChunks: true}),
    new DotenvPlugin(),
    new webpack.EnvironmentPlugin({NODE_ENV: "development"})
  ],
  output: {
    path: __dirname,
    filename: "bundle.js"
  },
  resolve: {
    extensions: [".js", ".jsx"]
  }
};
