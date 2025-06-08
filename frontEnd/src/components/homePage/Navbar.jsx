import { useNavigate } from "@solidjs/router";
import { createSignal, Switch, Match, Show } from "solid-js";
import { AiOutlineMenuFold, AiOutlineMenuUnfold } from "solid-icons/ai";

const Navbar = () => {
  const navigate = useNavigate();
  const [isMenuOpen, setIsMenuOpen] = createSignal(false);

  const scrollToSection = (sectionId) => {
    if (typeof window === "undefined") return;
    const element = document.getElementById(sectionId);
    if (element) {
      const headerElement = document.getElementById("header");
      const headerHeight = headerElement ? headerElement.offsetHeight : 0;
      const targetScrollPosition =
        element.getBoundingClientRect().top + window.scrollY - headerHeight;

      window.scrollTo({
        top: targetScrollPosition,
        behavior: "smooth",
      });
      setIsMenuOpen(false);
    }
  };

  return (
    <header
      id="header"
      className="sticky top-0 z-10 bg-gray-900/90 backdrop-blur-md py-4 px-6 border-b border-gray-800"
    >
      <div className="container mx-auto flex justify-between items-center">
        <div className="flex items-center">
          <img
            src="./BlueFoxLogo.png"
            alt="Blue Fox Logo"
            className="w-14 h-14 mr-2"
          />
          <span className="text-xl font-bold">
            Blue <span className="text-blue-400">Fox</span>
          </span>
        </div>

        {/* Mobile Menu Toggle */}
        <button
          className="md:hidden text-gray-300 hover:text-white transition-colors"
          onClick={() => setIsMenuOpen(!isMenuOpen())}
          aria-label="Toggle navigation menu"
        >
          <Switch>
            <Match when={isMenuOpen()}>
              <AiOutlineMenuUnfold className="w-6 h-6" />
            </Match>
            <Match when={!isMenuOpen()}>
              <AiOutlineMenuFold className="w-6 h-6" />
            </Match>
          </Switch>
        </button>

        {/* Desktop Navigation */}
        <nav className="hidden md:flex gap-6">
          <button
            onClick={() => scrollToSection("main")}
            className="text-gray-300 hover:text-white transition-colors"
            aria-label="Go to Main section"
          >
            Main
          </button>
          <button
            onClick={() => scrollToSection("features")}
            className="text-gray-300 hover:text-white transition-colors"
            aria-label="Go to Features section"
          >
            Features
          </button>
          <button
            onClick={() => scrollToSection("technologies-packages-hosting")}
            className="text-gray-300 hover:text-white transition-colors"
            aria-label="Go to Technologies section"
          >
            Technologies
          </button>
          <button
            onClick={() => scrollToSection("developer")}
            className="text-gray-300 hover:text-white transition-colors"
            aria-label="Go to Developer section"
          >
            Developer
          </button>
          <button
            onClick={() => scrollToSection("footer")}
            className="text-gray-300 hover:text-white transition-colors"
            aria-label="Go to Footer section"
          >
            Footer
          </button>
          <button
            size="sm"
            className="bg-gradient-to-r from-blue-500 to-blue-600 text-white hover:from-blue-600 hover:to-blue-700
                       px-4 py-2 rounded-full shadow-md hover:shadow-lg transition-all duration-300
                       font-semibold text-sm"
            onClick={() => navigate("/auth")}
          >
            Start Chatting
          </button>
        </nav>
      </div>

      {/* Mobile menu items */}
      <Show when={isMenuOpen()}>
        <nav className="md:hidden mt-4 space-y-2">
          <button
            onClick={() => scrollToSection("main")}
            className="block w-full text-left py-2 px-4 text-gray-300 hover:bg-gray-800 hover:text-white rounded-md transition-colors"
          >
            Main
          </button>
          <button
            onClick={() => scrollToSection("features")}
            className="block w-full text-left py-2 px-4 text-gray-300 hover:bg-gray-800 hover:text-white rounded-md transition-colors"
          >
            Features
          </button>
          <button
            onClick={() => scrollToSection("technologies-packages-hosting")}
            className="block w-full text-left py-2 px-4 text-gray-300 hover:bg-gray-800 hover:text-white rounded-md transition-colors"
          >
            Technologies
          </button>
          <button
            onClick={() => scrollToSection("developer")}
            className="block w-full text-left py-2 px-4 text-gray-300 hover:bg-gray-800 hover:text-white rounded-md transition-colors"
          >
            Developer
          </button>
          <button
            onClick={() => scrollToSection("footer")}
            className="block w-full text-left py-2 px-4 text-gray-300 hover:bg-gray-800 hover:text-white rounded-md transition-colors"
          >
            Footer
          </button>
          <button
            className="block w-full text-left py-2 px-4 bg-gradient-to-r from-blue-500 to-blue-600 text-white hover:from-blue-600 hover:to-blue-700 rounded-md shadow-md hover:shadow-lg transition-all duration-300 font-semibold text-sm"
            onClick={() => {
              navigate("/auth");
              setIsMenuOpen(false);
            }}
          >
            Start Chatting
          </button>
        </nav>
      </Show>
    </header>
  );
};

export default Navbar;
