import { Motion } from "solid-motionone";
import { createSignal, Switch, Match } from "solid-js";
import { onMount } from "solid-js";
import AuthHeader from "~/components/authPage/AuthHeader";
import AuthHomeButton from "~/components/authPage/AuthHomeButton";
import Input from "~/components/Input";
import AuthFooter from "~/components/authPage/AuthFooter";
import { TbCalendarCode, TbLock, TbMail, TbUser } from "solid-icons/tb";

const login = () => {
  const [email, setEmail] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [confirmPassword, setConfirmPassword] = createSignal("");
  const [username, setUsername] = createSignal("");
  const [birthDate, setbirthDate] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const [error, setError] = createSignal<string | null>(null);
  const [mounted, setMounted] = createSignal(false);

  const [view, setView] = createSignal<"login" | "signup" | "resetPassword">(
    "login"
  );

  onMount(() => {
    setMounted(true);
  });

  const handleAuth = async (e: SubmitEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    if (view() === "resetPassword") {
      console.log("Resetting password with:", {
        email: email(),
      });
      alert("Successfully reset password! (Simulated)");
      setLoading(false);
      return;
    }
    if (
      !email() ||
      !password() ||
      (view() === "signup" &&
        (!username() || !birthDate() || !confirmPassword()))
    ) {
      setError("Please fill in all fields.");
      setLoading(false);
      return;
    }

    if (view() === "signup" && password() !== confirmPassword()) {
      setError("Passwords do not match.");
      setLoading(false);
      return;
    }

    let formattedBirthDate = birthDate(); // Get the current value from the signal

    // Only attempt to format if birthDate() has a value
    if (formattedBirthDate) {
      try {
        // Create a Date object from the input string.
        // Using 'T00:00:00' ensures it's interpreted at the start of the day in local time
        // before converting to UTC.
        const dateObject = new Date(`${formattedBirthDate}T00:00:00`);

        // Check if the dateObject is valid
        if (!isNaN(dateObject.getTime())) {
          formattedBirthDate = dateObject.toISOString();
        } else {
          setError("Invalid Birth Date provided.");
          setLoading(false);
          return;
        }
      } catch (error) {
        setError("Error formatting Birth Date.");
        setLoading(false);
        return;
      }
    }

    try {
      const dateObject = new Date(`${birthDate()}T00:00:00`);
      if (view() === "signup") {
        const response = await fetch("http://localhost:9000/api/user", {
          method: "PUT",
          headers: {
            Accept: "application/json",
          },
          body: JSON.stringify({
            email: email(),
            password_hash: password(),
            username: username(),
            date_of_birth: formattedBirthDate,
          }),
        });
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
        <AuthHeader view={view} setView={setView} />
        <form onsubmit={(e) => handleAuth(e)} class="space-y-4">
          {
            <Switch>
              <Match when={view() === "signup"}>
                <>
                  <Input
                    mounted={mounted}
                    context={"Username"}
                    placeHolder={"johndoe"}
                    IconName={TbUser}
                    value={username}
                    setValue={setUsername}
                    type={"text"}
                    autocomplete={"username"}
                  />
                  <Input
                    mounted={mounted}
                    context={"Email"}
                    placeHolder={"johndoe@example.com"}
                    IconName={TbMail}
                    value={email}
                    setValue={setEmail}
                    type={"email"}
                    autocomplete={"email"}
                  />
                  <Input
                    mounted={mounted}
                    context={"Date of Birth"}
                    IconName={TbCalendarCode}
                    placeHolder={""}
                    value={birthDate}
                    setValue={setbirthDate}
                    type={"date"}
                    autocomplete={"bday"}
                  />
                  <Input
                    mounted={mounted}
                    context={"Password"}
                    IconName={TbLock}
                    placeHolder={"Pass123@"}
                    value={password}
                    setValue={setPassword}
                    type={"password"}
                    autocomplete={"new-password"}
                  />
                  <Input
                    mounted={mounted}
                    context={"Confirm Password"}
                    IconName={TbLock}
                    placeHolder={"Pass123@"}
                    value={confirmPassword}
                    setValue={setConfirmPassword}
                    type={"password"}
                    autocomplete={"new-password"}
                  />
                </>
              </Match>
              <Match when={view() === "resetPassword"}>
                <Input
                  mounted={mounted}
                  context={"Email"}
                  placeHolder={"johndoe@example.com"}
                  IconName={TbMail}
                  value={email}
                  setValue={setEmail}
                  type={"email"}
                  autocomplete={"email"}
                />
              </Match>
              <Match when={view() === "login"}>
                <Input
                  mounted={mounted}
                  context={"Email"}
                  placeHolder={"johndoe@example.com"}
                  IconName={TbMail}
                  value={email}
                  setValue={setEmail}
                  type={"email"}
                  autocomplete={"email"}
                />
                <Input
                  mounted={mounted}
                  context={"Password"}
                  IconName={TbLock}
                  placeHolder={"Pass123@"}
                  value={password}
                  setValue={setPassword}
                  type={"password"}
                  autocomplete={"current-password"}
                />
              </Match>
            </Switch>
          }

          <div class="text-right">
            <button
              onClick={() => {
                view() === "resetPassword"
                  ? setView("login")
                  : setView("resetPassword");
              }}
              class="text-blue-400 hover:underline text-sm transition-colors"
              type="button"
            >
              {view() === "resetPassword"
                ? "Back to Login"
                : "Forgot Password?"}
            </button>
          </div>

          <button
            type="submit"
            class="bg-gradient-to-r from-blue-500 to-blue-600 text-white hover:from-blue-600 hover:to-blue-700
                           px-6 py-3 rounded-full shadow-md hover:shadow-lg transition-all duration-300
                           font-semibold text-lg w-full hover:scale-105"
            disabled={loading()}
          >
            {
              <Switch>
                <Match when={view() === "signup"}>{"Sign Up"}</Match>
                <Match when={view() === "resetPassword"}>
                  {"Send Reset Link"}
                </Match>
                <Match when={view() === "login"}>{"Login"}</Match>
              </Switch>
            }
          </button>
        </form>

        <AuthFooter view={view} setView={setView} />
      </Motion.div>
    </div>
  );
};

export default login;
