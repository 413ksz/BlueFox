const AuthFooter = ({ isSignUp, setIsSignUp, setIsPasswordReset }) => {
  return (
    <div class="text-center text-gray-400 text-sm">
      {isSignUp() ? (
        <p>
          Already have an account?{" "}
          <button
            onClick={() => {
              setIsSignUp(false);
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
              setIsSignUp(true);
              setIsPasswordReset(false);
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
