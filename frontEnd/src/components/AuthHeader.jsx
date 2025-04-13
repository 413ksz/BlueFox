import { Motion } from "solid-motionone";
const AuthHeader = ({ isSignUp }) => {
  return (
    <div className="text-center mb-8">
      <Motion.h2
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: 20 }}
        transition={{ duration: 0.5, ease: "easeInOut" }}
        className="text-3xl font-bold text-white"
      >
        {isSignUp() ? "Join Us!" : "Welcome Back!"}
      </Motion.h2>
      <Motion.p
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: 10 }}
        transition={{ duration: 0.5, ease: "easeInOut", delay: 0.2 }}
        className="text-gray-300 text-sm"
      >
        {isSignUp()
          ? "Create an account to start your journey!"
          : "Log in to continue your adventure."}
      </Motion.p>
    </div>
  );
};

export default AuthHeader;
