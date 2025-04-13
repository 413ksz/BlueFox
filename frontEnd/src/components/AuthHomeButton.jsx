const AuthHomeButton = () => {
  return (
    <div className="absolute top-6 left-6">
      <button
        variant="ghost"
        onClick={() => (window.location.href = "/")}
        className="text-gray-300 hover:text-white transition-all duration-300
          flex items-center gap-2
          hover:bg-gray-800/50 hover:scale-105
          hover:shadow-lg"
      >
        <img src="./BlueFoxLogo.png" alt="Home" className="w-10 h-10" />
        <span className="hidden sm:inline">Home</span>
      </button>
    </div>
  );
};

export default AuthHomeButton;
