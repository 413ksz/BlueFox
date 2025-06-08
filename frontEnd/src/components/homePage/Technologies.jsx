import ItemMapper from "./ItemMapper";

const Technologies = () => {
  const technologies = [
    { name: "SolidJS", url: "https://www.solidjs.com/" },
    { name: "SolidStart", url: "https://start.solidjs.com/" },
    { name: "Tailwind CSS", url: "https://tailwindcss.com/" },
  ];
  const packages = [
    {
      name: "@solidjs/meta",
      url: "https://github.com/solidjs/solid-meta",
    },
    {
      name: "@solidjs/router",
      url: "https://github.com/solidjs/solid-router",
    },
    { name: "@solidjs/start", url: "https://start.solidjs.com/" },
    { name: "solid-icons", url: "https://solid-icons.vercel.app/" },
    {
      name: "solid-motionone",
      url: "https://github.com/solidjs-community/solid-motionone",
    },
    { name: "vinxi", url: "https://vinxi.vercel.app/" },
    {
      name: "@tailwindcss/postcss",
      url: "https://github.com/tailwindlabs/tailwindcss",
    },
    { name: "postcss", url: "https://postcss.org/" },
  ];

  const hostingPlatforms = [{ name: "Vercel", url: "https://vercel.com/" }];

  return (
    <section
      id="technologies-packages-hosting"
      className="bg-gray-950/50 py-16 md:py-24"
    >
      <div className="container mx-auto px-4 max-w-4xl">
        {/* Technologies Section */}
        <h2 className="text-3xl sm:text-4xl font-extrabold text-center mb-10 text-white leading-tight">
          Technologies & <span className="text-blue-400">Tools</span>
        </h2>
        <p className="text-lg text-gray-400 text-center mb-12 max-w-2xl mx-auto">
          The core frameworks and libraries that power the application.
        </p>

        <ItemMapper items={technologies} />

        {/* Packages Section */}
        <h2 className="text-3xl sm:text-4xl font-extrabold text-center mb-10 text-white leading-tight">
          Key <span className="text-blue-400">Packages</span>
        </h2>
        <p className="text-lg text-gray-400 text-center mb-12 max-w-2xl mx-auto">
          Essential utilities and specialized libraries enhancing functionality.
        </p>

        <ItemMapper items={packages} />

        {/* Hosting Platforms Section */}
        <h2 className="text-3xl sm:text-4xl font-extrabold text-center mb-10 text-white leading-tight">
          Hosting <span className="text-blue-400">Platforms</span>
        </h2>
        <p className="text-lg text-gray-400 text-center mb-12 max-w-2xl mx-auto">
          Where the application are deployed and hosted.
        </p>

        <ItemMapper items={hostingPlatforms} />
      </div>
    </section>
  );
};

export default Technologies;
