const NavLink = ({ scrollToSection, scrollToSectionId, label, isMobile }) => {
  return (
    <a
      href={`#${scrollToSectionId}`}
      onClick={(e) => {
        e.preventDefault(); // Stop the browser from jumping directly
        scrollToSection(scrollToSectionId);
      }}
      className={`text-gray-300 hover:text-white transition-colors ${
        isMobile
          ? "block w-full text-left py-2 px-4 hover:bg-gray-800 rounded-md"
          : "inline-block text-center px-1 py-1"
      }`}
      aria-label={`Go to ${label} section`}
    >
      {label}
    </a>
  );
};

export default NavLink;
