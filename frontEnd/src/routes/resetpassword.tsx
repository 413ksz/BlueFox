import { Motion } from "solid-motionone";
import { createSignal, onMount } from "solid-js";
import { IconLock } from "@tabler/icons-solidjs";
import AuthHomeButton from "~/components/authPage/AuthHomeButton";
import Input from "~/components/Input";
import { useNavigate } from "@solidjs/router";

const ResetPassword = () => {
  const [password, setPassword] = createSignal("");
  const [confirmPassword, setConfirmPassword] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const [error, setError] = createSignal<string | null>(null);
  const [mounted, setMounted] = createSignal(false);
  const navigate = useNavigate();

  onMount(() => {
    setMounted(true);
  });

  const handleReset = async (e: SubmitEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    if (!password() || !confirmPassword()) {
      setError("Please fill in both password fields.");
      setLoading(false);
      return;
    }

    if (password() !== confirmPassword()) {
      setError("Passwords do not match.");
      setLoading(false);
      return;
    }

    try {
      // Simulate sending a reset password request
      await new Promise((resolve) => setTimeout(resolve, 1500));
      console.log("Password reset requested with:", { password: password() });
      alert("Password reset successful! (Simulated)");
      navigate("/auth");
    } catch (err: any) {
      setError(
        err.message || "An error occurred while resetting your password."
      );
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
      >
        <form onsubmit={(e) => handleReset(e)} class="space-y-4">
          <Input
            mounted={mounted}
            context={"New Password"}
            placeHolder={"NewPass123@"}
            IconName={IconLock}
            value={password}
            setValue={setPassword}
            type={"password"}
            autocomplete={"new-password"}
          />
          <Input
            mounted={mounted}
            context={"Confirm New Password"}
            placeHolder={"NewPass123@"}
            IconName={IconLock}
            value={confirmPassword}
            setValue={setConfirmPassword}
            type={"password"}
            autocomplete={"new-password"}
          />

          {error() && <p class="text-red-400 text-sm">{error()}</p>}

          <button
            type="submit"
            class="bg-gradient-to-r from-blue-500 to-blue-600 text-white hover:from-blue-600 hover:to-blue-700
                   px-6 py-3 rounded-full shadow-md hover:shadow-lg transition-all duration-300
                   font-semibold text-lg w-full hover:scale-105"
            disabled={loading()}
          >
            {loading() ? "Resetting..." : "Reset Password"}
          </button>
        </form>
      </Motion.div>
    </div>
  );
};

export default ResetPassword;
