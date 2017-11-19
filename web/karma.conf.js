// Karma configuration
var webpack = require("webpack");
var DotenvPlugin = require("webpack-dotenv-plugin");

if (!process.env.NODE_ENV) {
  process.env.NODE_ENV = "test";
}

module.exports = function(config) {
  config.set({
    // base path that will be used to resolve all patterns (eg. files, exclude)
    // basePath: "",
    // frameworks to use
    // available frameworks: https://npmjs.org/browse/keyword/karma-adapter
    frameworks: ["mocha"],
    // list of files / patterns to load in the browser
    files: ["test.webpack.js"],
    // list of files to exclude
    // exclude: [],
    // preprocess matching files before serving them to the browser
    // available preprocessors: https://npmjs.org/browse/keyword/karma-preprocessor
    preprocessors: {
      "test.webpack.js": ["webpack", "sourcemap"]
    },
    // test results reporter to use
    // possible values: 'dots', 'progress'
    // available reporters: https://npmjs.org/browse/keyword/karma-reporter
    reporters: ["mocha", "notify", "coverage", "clear-screen"],
    mochaReporter: {
      showDiff: true
    },

    // web server port
    port: 9876,
    // enable / disable colors in the output (reporters and logs)
    colors: true,
    // level of logging
    // possible values: config.LOG_DISABLE || config.LOG_ERROR || config.LOG_WARN || config.LOG_INFO || config.LOG_DEBUG
    logLevel: config.LOG_INFO,
    browserConsoleLogOptions: {
      level: "log",
      format: "%b %T: %m",
      terminal: true
    },

    // enable / disable watching file and executing tests whenever any file changes
    autoWatch: true,
    // start these browsers
    // available browser launchers: https://npmjs.org/browse/keyword/karma-launcher
    browsers: ["Chrome"],
    webpack: {
      devtool: "inline-source-map",
      externals: {
        // Enzyme includes require statements for these modules for backwards compatibility
        // with older versions of React. Webpack gets confused by these, even though
        // they will never actually be required. We are Marking them as externals
        // so webpack doesn't complain.
        "react/addons": true,
        "react/lib/ExecutionEnvironment": true,
        "react/lib/ReactContext": true
      },
      module: {
        // rename to rules
        rules: [
          // {
          //   test: /\.json$/,
          //   loader: "json-loader",
          // },
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
            use: [
              {loader: "style-loader"},
              {loader: "css-loader"},
              {loader: "sass-loader"}
            ]
          }
        ]
      },
      resolve: {
        extensions: [".js", ".jsx", ".json"]
      },
      plugins: [new DotenvPlugin()],
      watch: true
    },
    webpackServer: {
      noInfo: true
    },

    // Continuous Integration mode
    // if true, Karma captures browsers, runs the tests and exits
    singleRun: false,

    // Concurrency level
    // how many browser should be started simultaneous
    concurrency: Infinity,
    coverageReporter: {
      type: "html", //produces a html document after code is run
      dir: "coverage" //path to created html doc
    }
  });
};
