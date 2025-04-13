import { createSignal, onMount } from "solid-js";
import { Motion } from "solid-motionone";

function Hero() {
  const [mounted, setMounted] = createSignal(false);

  onMount(() => {
    setMounted(true);
  });

  const scrollToSection = (id) => {
    const element = document.getElementById(id);
    element?.scrollIntoView({ behavior: "smooth" });
  };

  return (
    <main
      id="main"
      className="flex-1 flex items-center justify-center p-4 md:p-8"
    >
      <div className="text-center space-y-6 max-w-3xl">
        <Motion.h1
          initial={{ opacity: 0, y: -20 }}
          animate={mounted() ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.6, easing: "ease-in-out" }}
          className="text-4xl sm:text-5xl md:text-6xl font-extrabold bg-clip-text text-transparent bg-gradient-to-r from-blue-500 to-cyan-500"
        >
          Welcome to Blue Fox
        </Motion.h1>
        <Motion.p
          initial={{ opacity: 0, y: 20 }}
          animate={mounted() ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.6, easing: "ease-in-out", delay: 0.2 }}
          className="text-lg sm:text-xl text-gray-200"
        >
          A modern, open-source chat application built with SolidJS and Tailwind
          CSS.
        </Motion.p>
        <Motion.div
          initial={{ opacity: 0 }}
          animate={mounted() ? { opacity: 1 } : {}}
          transition={{ duration: 0.6, easing: "ease-in-out", delay: 0.4 }}
        >
          <button
            size="lg"
            className="bg-gradient-to-r from-blue-500 to-cyan-500 text-white hover:from-blue-600 hover:to-cyan-600
                     px-8 py-3 rounded-full shadow-lg hover:shadow-xl transition-all duration-300
                     font-semibold text-lg flex items-center gap-2"
            onClick={() => scrollToSection("features")}
          >
            Get Started
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="w-5 h-5"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M8.25 4.5l7.5 7.5-7.5 7.5"
              />
            </svg>
          </button>
        </Motion.div>
      </div>
    </main>
  );
}

export default Hero;
