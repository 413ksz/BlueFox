const Navbar = () => {
  const scrollToSection = (sectionId) => {
    if (typeof window === "undefined") return;
    const element = document.getElementById(sectionId);
    if (element) {
      element.scrollIntoView({ behavior: "smooth" });
    }
  };
  return (
    <header className="sticky top-0 z-10 bg-gray-900/90 backdrop-blur-md py-4 px-6 border-b border-gray-800">
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
        <nav className="flex gap-6">
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
            onClick={() => scrollToSection("technologies")}
            className="text-gray-300 hover:text-white transition-colors"
            aria-label="Go to Features section"
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
            onClick={() => (window.location.href = "/auth")}
          >
            Start Chatting
          </button>
        </nav>
      </div>
    </header>
  );
};

export default Navbar;
