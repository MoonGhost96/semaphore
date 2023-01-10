module.exports = {
  configureWebpack: {
    devServer: {
      historyApiFallback: true,
      // when in dev env:
      proxy: {
        '^/api': {
          target: 'http://localhost:3000/',
          ws: true,
          changeOrigin: true,
        },
      },
    },
  },
  chainWebpack: (config) => {
    config.plugin('html')
      .tap((args) => {
        // eslint-disable-next-line no-param-reassign
        args[0].minify = false;
        return args;
      });
  },
  transpileDependencies: [
    'vuetify',
  ],
  publicPath: './',
};
