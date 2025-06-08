const NavLink = ({ scrollToSection, scrollToSectionId, label, isMobile }) => {
  return (
    <button
      onClick={() => scrollToSection(scrollToSectionId)}
      className={`text-gray-300 hover:text-white transition-colors ${
        isMobile
          ? "block w-full text-left py-2 px-4 hover:bg-gray-800 rounded-md"
          : ""
      }`}
      aria-label={`Go to ${scrollToSectionId} section`}
    >
      {label}
    </button>
  );
};

export default NavLink;
