import nextra from 'nextra';

const withNextra = nextra({
  theme: 'nextra-theme-docs',
  themeConfig: './theme.config.jsx',
  defaultShowCopyCode: true,
});

export default withNextra({
  reactStrictMode: true,
  images: {
    unoptimized: true,
  },
  output: 'export',
  distDir: 'dist',
  basePath: '/' + process.env.NEXT_PUBLIC_REPOSITORY.split('/')[1],
});
