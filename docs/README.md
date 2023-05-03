# Documentation

This folder contains the documentation.

## Build logo

You can use this prompt to generate a logo on MideJourney for example:
`saas logo, assistant octopus, 2d very minimalist illustration, glitchy, cryengine, ironical, mote kei, sketchfab, blue and purple gradient --v 5`

## Development

Below a guide to help you run the website locally and preview changes before submitting them.

Please ensure that your environment match the following requirements :

- [pnpm](https://pnpm.io/installation)
- [NodeJS](https://nodejs.org/en/) (v17+)

To run the website locally :

```bash
export NEXT_PUBLIC_REPOSITORY=leofvo/gosurp
pnpm install
pnpm dev
```

The website should now be up and running at http://localhost:3000/gosurp. Next.js will automatically reload your changes. When you're done, stop the server by hitting `CTRL+C`.

## Setup

Change config in `next.config.js` to match your needs.
