// This file is the main entry point for the Karma Webpack test harness.
// It requires all the .jsx files in the test folder.

var context = require.context("./test/", true, /^.*\.jsx?$/);
context.keys().forEach(context);
