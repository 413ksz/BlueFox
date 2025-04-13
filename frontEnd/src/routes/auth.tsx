import { Motion } from "solid-motionone";
import { createSignal } from "solid-js";
import {
  IconUser,
  IconMail,
  IconCalendarClock,
  IconLock,
} from "@tabler/icons-solidjs";
import { onMount } from "solid-js";
import AuthHeader from "~/components/AuthHeader";
import AuthHomeButton from "~/components/AuthHomeButton";
import Input from "~/components/Input";
import AuthFooter from "~/components/AuthFooter";
const login = () => {
  const [email, setEmail] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [confirmPassword, setConfirmPassword] = createSignal("");
  const [name, setName] = createSignal("");
  const [birthDate, setbirthDate] = createSignal("");
  const [isSignUp, setIsSignUp] = createSignal(false);
  const [loading, setLoading] = createSignal(false);
  const [error, setError] = createSignal<string | null>(null);
  const [mounted, setMounted] = createSignal(false);

  onMount(() => {
    setMounted(true);
  });

  const handleAuth = async (e: SubmitEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    if (
      !email() ||
      !password() ||
      (isSignUp() && (!name() || !birthDate() || !confirmPassword()))
    ) {
      setError("Please fill in all fields.");
      setLoading(false);
      return;
    }

    if (isSignUp() && password() !== confirmPassword()) {
      setError("Passwords do not match.");
      setLoading(false);
      return;
    }

    try {
      await new Promise((resolve) => setTimeout(resolve, 1500));
      if (isSignUp()) {
        console.log("Signing up with:", {
          mame: name(),
          email: email(),
          birthDate: birthDate(),
          password: password(),
        });
        alert("Successfully signed up! (Simulated)");
      } else {
        console.log("Logging in with:", {
          email: email(),
          password: password(),
        });
        alert("Successfully logged in! (Simulated)");
      }
    } catch (err: any) {
      setError(err.message || "An error occurred.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div class="min-h-screen bg-gradient-to-br from-gray-900 via-blue-900 to-blue-950 text-white flex flex-col justify-center items-center p-4">
      <AuthHomeButton />
      <Motion.div
        class="w-full max-w-md bg-gray-900/90 backdrop-blur-md rounded-xl shadow-2xl p-6 md:p-8 border border-gray-800 space-y-6"
        initial={{ opacity: 0 }}
        animate={{ opacity: mounted() ? 1 : 0 }}
        transition={{ delay: 0.3 }}
      >
        <AuthHeader isSignUp={isSignUp} />
        <form onsubmit={(e) => handleAuth(e)} class="space-y-4">
          {isSignUp() ? (
            <>
              <Input
                mounted={mounted}
                context={"Name"}
                placeHolder={"John Doe"}
                IconName={IconUser}
                value={name}
                setValue={setName}
                type={"text"}
              />
              <Input
                mounted={mounted}
                context={"Email"}
                placeHolder={"johndoe@example.com"}
                IconName={IconMail}
                value={email}
                setValue={setEmail}
                type={"email"}
              />
              <Input
                mounted={mounted}
                context={"Date of Birth"}
                IconName={IconCalendarClock}
                placeHolder={""}
                value={birthDate}
                setValue={setbirthDate}
                type={"date"}
              />
              <Input
                mounted={mounted}
                context={"Password"}
                IconName={IconLock}
                placeHolder={"Pass123@"}
                value={password}
                setValue={setPassword}
                type={"password"}
              />
              <Input
                mounted={mounted}
                context={"Confirm Password"}
                IconName={IconLock}
                placeHolder={"Pass123@"}
                value={confirmPassword}
                setValue={setConfirmPassword}
                type={"password"}
              />
            </>
          ) : (
            <>
              <Input
                mounted={mounted}
                context={"Email"}
                placeHolder={"johndoe@example.com"}
                IconName={IconMail}
                value={email}
                setValue={setEmail}
                type={"email"}
              />
              <Input
                mounted={mounted}
                context={"Password"}
                IconName={IconLock}
                placeHolder={"Pass123@"}
                value={password}
                setValue={setPassword}
                type={"password"}
              />
              <div class="text-right">
                <button
                  onClick={() => {
                    alert("Forgot Password (Simulated)");
                  }}
                  class="text-blue-400 hover:underline text-sm transition-colors"
                >
                  Forgot Password?
                </button>
              </div>
            </>
          )}

          <button
            type="submit"
            class="bg-gradient-to-r from-blue-500 to-blue-600 text-white hover:from-blue-600 hover:to-blue-700 
                                     px-6 py-3 rounded-full shadow-md hover:shadow-lg transition-all duration-300 
                                     font-semibold text-lg w-full hover:scale-105"
            disabled={loading()}
          >
            {isSignUp() ? "Sign Up" : "Login"}
          </button>
        </form>

        <AuthFooter isSignUp={isSignUp} setIsSignUp={setIsSignUp} />
      </Motion.div>
    </div>
  );
};

export default login;
