# SolidStart Project

"This is a [SolidStart](https://docs.solidjs.com/solid-start) project, a robust framework that gives you everything you need to build powerful and performant Solid applications. Known for its fine-grained reactivity and minimal overhead, SolidStart streamlines your development workflow, providing features like server-side rendering, routing, and more."

# Getting Started
To get started with this project, follow the steps below.



## Creating a New Project
If you're setting up a new SolidStart project from scratch, you can use one of the following commands. This will create a new directory (or use the current one) and set up the basic project structure.

```bash
# Create a new project in the current directory
npm init solid@latest
bun create solid

# Create a new project in a specified directory (e.g., 'my-app')
npm init solid@latest my-app
bun create solid my-app
```

## Installation

Once you've created or cloned a SolidStart project, navigate into the project directory and install the necessary dependencies using your preferred package manager:

```bash
npm install
# or
pnpm install
# or
yarn install
# or
bun install
```
# Development
After installing dependencies, you can start the development server. This will compile your application and make it accessible in your browser, with hot-reloading for a smooth development experience.
```bash
npm run dev
# or
bun run dev

# To start the server and automatically open the app in a new browser tab:
npm run dev -- --open
```


# Building for Production
SolidStart applications are built using presets, which optimize your project for deployment to various environments (e.g., Node.js, Vercel, Netlify).
By default, running npm run build will generate a Node.js application that you can run with npm start.
To deploy to a different environment, you'll need to:
1. Install the relevant preset: Add it to your devDependencies in package.json (e.g., @solidjs/start-vercel).
2. Configure app.config.js: Specify the chosen preset in your app.config.js file.

##Built with
This project was initially set up using the official [Solid CLI](https://solid-cli.netlify.app), a command-line interface designed to help you quickly scaffold and manage Solid projects.
