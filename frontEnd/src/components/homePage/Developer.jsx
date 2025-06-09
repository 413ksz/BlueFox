import { TbBrandGithub } from "solid-icons/tb";
import { createSignal, onMount } from "solid-js";
import { Motion } from "solid-motionone";
const Developer = () => {
  const [mounted, setMounted] = createSignal(false);

  onMount(() => {
    setMounted(true);
  });
  const developer = {
    name: "413ksz",
    avatar: "https://github.com/413ksz.png",
    bio: "A full-stack developer and creator of Blue Fox. Loves building performant and user-friendly applications.",
    github: "https://github.com/413ksz",
  };
  return (
    <section id="developer" className="bg-gray-950/50 py-16 md:py-24">
      <div className="container mx-auto px-4">
        <Motion.h2
          initial={{ opacity: 0, y: -20 }}
          animate={mounted() ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.6, ease: "easeOutCubic" }}
          className="text-3xl sm:text-4xl font-semibold text-center mb-12 text-white"
        >
          About the Developer
        </Motion.h2>
        <Motion.div
          initial={{ opacity: 0, y: 50 }}
          animate={mounted() ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.7, ease: "easeOutCubic", delay: 0.2 }}
          className="flex flex-col md:flex-row items-center gap-8 bg-gray-800/50 rounded-xl p-6 shadow-lg border border-gray-700"
        >
          <Motion.img
            initial={{ opacity: 0 }}
            animate={mounted() ? { opacity: 1 } : {}}
            transition={{ duration: 0.5, delay: 0.4 }}
            src={developer.avatar}
            alt={developer.name}
            className="rounded-full w-40 h-40 border-4 border-gray-700 shadow-md transition-all duration-300 hover:scale-105"
          />
          <div className="text-center md:text-left space-y-4">
            <Motion.h3
              initial={{ opacity: 0, x: -20 }}
              animate={mounted() ? { opacity: 1, x: 0 } : {}}
              transition={{ duration: 0.5, delay: 0.5 }}
              className="text-2xl font-semibold text-white"
            >
              {developer.name}
            </Motion.h3>
            <Motion.p
              initial={{ opacity: 0, x: -20 }}
              animate={mounted() ? { opacity: 1, x: 0 } : {}}
              transition={{ duration: 0.5, delay: 0.6 }}
              className="text-gray-300 max-w-xl"
            >
              {developer.bio}
            </Motion.p>
            <Motion.div
              initial={{ opacity: 0 }}
              animate={mounted() ? { opacity: 1 } : {}}
              transition={{
                duration: 0.4,
                delay: 0.7,
                ease: "linear",
              }}
              className="flex justify-center md:justify-start gap-4"
            >
              <a
                href={developer.github}
                target="_blank"
                rel="noopener noreferrer"
                className="hover:text-blue-400 transition-colors flex items-center gap-1"
                aria-label="GitHub Profile"
              >
                {mounted() && <TbBrandGithub className="w-6 h-6" />}
                <span className="font-medium">GitHub</span>
              </a>
            </Motion.div>
          </div>
        </Motion.div>
      </div>
    </section>
  );
};
export default Developer;
