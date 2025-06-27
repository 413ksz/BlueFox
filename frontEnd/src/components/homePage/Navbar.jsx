import { A } from "@solidjs/router";
import { createSignal, Switch, Match, Show, For } from "solid-js";
import { AiOutlineMenuFold, AiOutlineMenuUnfold } from "solid-icons/ai";
import NavLink from "./NavLink";

const Navbar = () => {
  const [isMenuOpen, setIsMenuOpen] = createSignal(false);

  const navItems = [
    { id: "main", label: "Main" },
    { id: "features", label: "Features" },
    { id: "technologies-packages-hosting", label: "Technologies" },
    { id: "developer", label: "Developer" },
    { id: "footer", label: "Footer" },
  ];

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
            src="./BlueFoxLogo.webp"
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
          aria-expanded={isMenuOpen()}
          aria-controls="mobile-menu"
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
          <For each={navItems}>
            {(item) => (
              <NavLink
                scrollToSection={scrollToSection}
                scrollToSectionId={item.id}
                label={item.label}
                isMobile={false}
              />
            )}
          </For>
          <A
            href="/auth"
            onClick={() => {
              setIsMenuOpen(false);
            }}
            size="sm"
            className="bg-gradient-to-r from-blue-500 to-blue-600 text-white hover:from-blue-600 hover:to-blue-700
                       px-4 py-2 rounded-full shadow-md hover:shadow-lg transition-all duration-300
                       font-semibold text-sm"
          >
            Start Chatting
          </A>
        </nav>
      </div>

      {/* Mobile menu items */}
      <Show when={isMenuOpen()}>
        <nav id="mobile-menu" className="md:hidden mt-4 space-y-2">
          <For each={navItems}>
            {(item) => (
              <NavLink
                scrollToSection={scrollToSection}
                scrollToSectionId={item.id}
                label={item.label}
                isMobile={true}
              />
            )}
          </For>
          <A
            href="/auth"
            className="block w-full text-left py-2 px-4 bg-gradient-to-r from-blue-500 to-blue-600 text-white hover:from-blue-600 hover:to-blue-700 rounded-md shadow-md hover:shadow-lg transition-all duration-300 font-semibold text-sm"
            onClick={() => {
              setIsMenuOpen(false);
            }}
          >
            Start Chatting
          </A>
        </nav>
      </Show>
    </header>
  );
};

export default Navbar;
