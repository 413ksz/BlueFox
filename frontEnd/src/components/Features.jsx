import { Motion } from "solid-motionone";
import {
  IconBolt,
  IconSend,
  IconBrandGithub,
  IconUsers,
  IconRocket,
  IconBook,
} from "@tabler/icons-solidjs";
import { createSignal, onMount } from "solid-js";

const Features = () => {
      const [mounted, setMounted] = createSignal(false);
    
      onMount(() => {
        setMounted(true);
      });
  const features = [
    {
      title: "Blazing Fast",
      description: "Experience lightning-fast performance with SolidJS.",
      icon: IconBolt,
    },
    {
      title: "Real-time Chat",
      description:
        "Connect instantly with our efficient, low-latency chat engine.",
      icon: IconSend,
    },
    {
      title: "Open Source",
      description:
        "Explore and contribute to this project on GitHub.",
      icon: IconBrandGithub,
    },
    {
      title: "Community Driven",
      description: "Join a vibrant community of developers users.",
      icon: IconUsers,
    },
    {
      title: "Lightweight",
      description:
        "Minimal footprint for maximum performance and faster load times.",
      icon: IconRocket,
    },
    {
      title: "Documentation",
      description:
        "Comprehensive documentation to help you get started quickly.",
      icon: IconBook,
    },
  ];
  return (
    <main id="features" className="bg-gray-950/50 py-16 md:py-24 w-full">
    <div className="container mx-auto px-4">
        <h2 className="text-3xl sm:text-4xl font-semibold text-center mb-12 text-white">Key Features</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8">
            {features.map((feature, index) => (
              <Motion.div
              key={index}
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ duration: 0.4, delay: index * 0.08, ease: "linear" }}
              className={"p-6 rounded-xl border border-gray-800 shadow-lg bg-gray-900/80 backdrop-blur-md hover:border-blue-500/30 hover:shadow-blue-500/20 transition-all duration-300 flex flex-col items-center text-center space-y-4 group"}
            >
                    {
                        mounted() &&  <feature.icon className="w-10 h-10 text-blue-400 transition-transform group-hover:scale-110" aria-hidden="true" />
                    }
                    <h3 className="text-xl font-semibold text-white">{feature.title}</h3>
                    <p className="text-gray-300">{feature.description}</p>
                </Motion.div>
            ))}
        </div>
    </div>
</main>
  );
};

export default Features;
