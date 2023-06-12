function Logo() {
  return (
    <div
      style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
      }}
    >
      <img
        style={{ borderRadius: '5px', marginRight: '10px' }}
        width={50}
        src={
          '/' + process.env.NEXT_PUBLIC_REPOSITORY.split('/')[1] + '/logo.png'
        }
      ></img>
      <h1 style={{ fontWeight: 'bold' }}>
        {process.env.NEXT_PUBLIC_REPOSITORY.split('/')[1]}
      </h1>
    </div>
  );
}

export default {
  logo: <Logo />,
  primaryHue: {
    dark: 199,
    light: 199,
  },
  docsRepositoryBase: `https://github.com/${process.env.NEXT_PUBLIC_REPOSITORY}/tree/main/docs`,
  project: {
    link: `https://github.com/${process.env.NEXT_PUBLIC_REPOSITORY}`,
  },
  editLink: {
    text: 'Edit this page on GitHub →',
  },
  footer: {
    text: `MIT 2023 © ${process.env.NEXT_PUBLIC_REPOSITORY.split('/')[0]}`,
  },
  useNextSeoProps() {
    return {
      titleTemplate: '%s – ' + process.env.NEXT_PUBLIC_REPOSITORY.split('/')[1],
    };
  },
  banner: {
    key: 'first-release',
    text: (
      <a
        href={'https://github.com/' + process.env.NEXT_PUBLIC_REPOSITORY}
        target="_blank"
      >
        🎉 {process.env.NEXT_PUBLIC_REPOSITORY} first version is released. Read
        more →
      </a>
    ),
  },
};
