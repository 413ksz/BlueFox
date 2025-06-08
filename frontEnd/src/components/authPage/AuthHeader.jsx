import { Match, Switch } from "solid-js";
import { Motion } from "solid-motionone";
const AuthHeader = ({ view, setView }) => {
  return (
    <div className="text-center mb-8">
      <Motion.h2
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: 20 }}
        transition={{ duration: 0.5, ease: "easeInOut" }}
        className="text-3xl font-bold text-white"
      >
        {
          <Switch>
            <Match when={view() === "signup"}>Join Us!</Match>
            <Match when={view() === "login"}>Welcome Back!</Match>
            <Match when={view() === "resetPassword"}>
              Let's get you back in!
            </Match>
          </Switch>
        }
      </Motion.h2>
      <Motion.p
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: 10 }}
        transition={{ duration: 0.5, ease: "easeInOut", delay: 0.2 }}
        className="text-gray-300 text-sm"
      >
        {
          <Switch>
            <Match when={view() === "signup"}>
              Create an account to start your journey!
            </Match>
            <Match when={view() === "login"}>
              Log in to continue your adventure!
            </Match>
            <Match when={view() === "resetPassword"}>
              Enter the email you use for your account, and we'll send a reset
              link.
            </Match>
          </Switch>
        }
      </Motion.p>
    </div>
  );
};

export default AuthHeader;
