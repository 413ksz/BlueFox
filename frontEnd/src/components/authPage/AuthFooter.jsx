const AuthFooter = ({ view, setView }) => {
  return (
    <div class="text-center text-gray-400 text-sm">
      {view() === "signup" ? (
        <p>
          Already have an account?{" "}
          <button
            onClick={() => {
              setView("login");
            }}
            class="text-blue-400 hover:underline font-medium transition-colors"
          >
            Login
          </button>
        </p>
      ) : (
        <p>
          Don't have an account?{" "}
          <button
            onClick={() => {
              setView("signup");
            }}
            class="text-blue-400 hover:underline font-medium transition-colors"
          >
            Sign Up
          </button>
        </p>
      )}
    </div>
  );
};

export default AuthFooter;
