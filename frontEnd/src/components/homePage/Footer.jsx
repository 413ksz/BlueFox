import { TbBrandGithub } from "solid-icons/tb";
import { createSignal, onMount } from "solid-js";
const Footer = () => {
  const [mounted, setMounted] = createSignal(false);

  onMount(() => {
    setMounted(true);
  });
  return (
    <footer
      id="footer"
      className="py-8 text-center text-gray-400 border-t border-gray-800"
    >
      <div className="container mx-auto px-4">
        <p>&copy; {new Date().getFullYear()} Blue Fox. All rights reserved.</p>
        <div className="mt-4 flex justify-center gap-4">
          <a
            href="https://github.com/413ksz/BlueFox"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-blue-400 transition-colors"
            aria-label="GitHub Repository"
          >
            {mounted() && <TbBrandGithub className="w-6 h-6" />}
          </a>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
