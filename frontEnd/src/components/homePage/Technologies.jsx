import ItemMapper from "./ItemMapper";
import { TbBrandSolidjs } from "solid-icons/tb";
import { VsServerEnvironment } from "solid-icons/vs";
import { SiTailwindcss, SiPostcss, SiVercel } from "solid-icons/si";
import { createSignal, onMount } from "solid-js";
import { Motion } from "solid-motionone";

const Technologies = () => {
  const technologies = [
    {
      name: "SolidJS",
      url: "https://www.solidjs.com/",
      icon: TbBrandSolidjs,
    },
    {
      name: "SolidStart",
      url: "https://start.solidjs.com/",
      icon: TbBrandSolidjs,
    },
    {
      name: "Tailwind CSS",
      url: "https://tailwindcss.com/",
      icon: SiTailwindcss,
    },
  ];
  const packages = [
    {
      name: "@solidjs/meta",
      url: "https://github.com/solidjs/solid-meta",
      icon: TbBrandSolidjs,
    },
    {
      name: "@solidjs/router",
      url: "https://github.com/solidjs/solid-router",
      icon: TbBrandSolidjs,
    },
    {
      name: "@solidjs/start",
      url: "https://start.solidjs.com/",
      icon: TbBrandSolidjs,
    },
    {
      name: "solid-icons",
      url: "https://solid-icons.vercel.app/",
      icon: TbBrandSolidjs,
    },

    {
      name: "solid-motionone",
      url: "https://github.com/solidjs-community/solid-motionone",
      icon: TbBrandSolidjs,
    },
    {
      name: "vinxi",
      url: "https://vinxi.vercel.app/",
      icon: VsServerEnvironment,
    },
    {
      name: "@tailwindcss/postcss",
      url: "https://github.com/tailwindlabs/tailwindcss",
      icon: SiTailwindcss,
    },
    { name: "postcss", url: "https://postcss.org/", icon: SiPostcss },
  ];

  const hostingPlatforms = [
    { name: "Vercel", url: "https://vercel.com/", icon: SiVercel },
  ];

  const [mounted, setMounted] = createSignal(false);

  onMount(() => {
    setMounted(true);
  });

  return (
    <section
      id="technologies-packages-hosting"
      className="bg-gray-950/50 py-16 md:py-24"
    >
      <div className="container mx-auto px-4 max-w-4xl">
        <Motion.h2
          initial={{ opacity: 0, y: -20 }}
          animate={mounted() ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.6, ease: "easeOutCubic" }}
          className="text-3xl sm:text-4xl font-extrabold text-center mb-10 text-white leading-tight"
        >
          Technologies & <span className="text-blue-400">Tools</span>
        </Motion.h2>
        <Motion.p
          initial={{ opacity: 0, y: -20 }}
          animate={mounted() ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.6, ease: "easeOutCubic", delay: 0.1 }}
          className="text-lg text-gray-400 text-center mb-12 max-w-2xl mx-auto"
        >
          The core frameworks and libraries that power the application.
        </Motion.p>

        <ItemMapper items={technologies} mounted={mounted} />

        <Motion.h2
          initial={{ opacity: 0, y: -20 }}
          animate={mounted() ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.6, ease: "easeOutCubic", delay: 0.2 }}
          className="text-3xl sm:text-4xl font-extrabold text-center mb-10 text-white leading-tight mt-16"
        >
          Key <span className="text-blue-400">Packages</span>
        </Motion.h2>
        <Motion.p
          initial={{ opacity: 0, y: -20 }}
          animate={mounted() ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.6, ease: "easeOutCubic", delay: 0.3 }}
          className="text-lg text-gray-400 text-center mb-12 max-w-2xl mx-auto"
        >
          Essential utilities and specialized libraries enhancing functionality.
        </Motion.p>

        <ItemMapper items={packages} mounted={mounted} />

        <Motion.h2
          initial={{ opacity: 0, y: -20 }}
          animate={mounted() ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.6, ease: "easeOutCubic", delay: 0.4 }}
          className="text-3xl sm:text-4xl font-extrabold text-center mb-10 text-white leading-tight mt-16"
        >
          Hosting <span className="text-blue-400">Platforms</span>
        </Motion.h2>
        <Motion.p
          initial={{ opacity: 0, y: -20 }}
          animate={mounted() ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.6, ease: "easeOutCubic", delay: 0.5 }}
          className="text-lg text-gray-400 text-center mb-12 max-w-2xl mx-auto"
        >
          Where the application are deployed and hosted.
        </Motion.p>

        <ItemMapper items={hostingPlatforms} mounted={mounted} />
      </div>
    </section>
  );
};

export default Technologies;
