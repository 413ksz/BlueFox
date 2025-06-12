import { createSignal, Show } from "solid-js";
import { Title } from "@solidjs/meta";

export default function App() {
  const [responseData, setResponseData] = createSignal(null);
  const [loading, setLoading] = createSignal(false);
  const [error, setError] = createSignal(null);

  const fetchData = async () => {
    setLoading(true);
    setError(null);
    setResponseData(null);

    try {
      const response = await fetch("/api/test", {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
      });

      console.log("Response:", response);

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(
          `HTTP error! Status: ${response.status} - ${errorText}`
        );
      }

      const data = await response.json();
      setResponseData(data);
    } catch (err: any) {
      console.error("Failed to fetch data:", err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div class="min-h-screen bg-gray-950 text-gray-100 font-inter flex flex-col items-center justify-center p-4">
      <Title>API Test Endpoint</Title>
      <div class="w-full max-w-2xl bg-gray-900 p-8 rounded-xl shadow-lg border border-gray-800 flex flex-col items-center gap-8">
        <h1 class="text-4xl md:text-5xl font-extrabold text-white text-center tracking-tight leading-tight">
          API Test Endpoint
        </h1>

        <p class="text-gray-400 text-lg text-center max-w-md">
          Click the button below to fetch data from the `/api/test` API route.
        </p>

        <button
          onClick={fetchData}
          disabled={loading()}
          class="px-8 py-4 bg-gradient-to-r from-blue-600 to-purple-600 text-white font-semibold rounded-full shadow-lg
                 hover:from-blue-700 hover:to-purple-700 transition-all duration-300 transform hover:-translate-y-1
                 focus:outline-none focus:ring-4 focus:ring-blue-500 focus:ring-opacity-50 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading() ? "Fetching..." : "Fetch API Data"}
        </button>

        <Show when={loading()}>
          <div class="text-blue-400 text-lg mt-4 animate-pulse">
            Loading data...
          </div>
        </Show>

        <Show when={error()}>
          <div class="bg-red-900 border border-red-700 text-red-200 p-4 rounded-lg w-full mt-4 text-center break-words">
            Error: {error()}
          </div>
        </Show>

        <Show when={responseData()}>
          <div class="bg-gray-800 p-6 rounded-lg w-full mt-4 shadow-inner flex flex-col items-center">
            <h2 class="text-2xl font-bold text-gray-200 mb-4">API Response:</h2>
            <pre class="bg-gray-900 text-gray-100 p-4 rounded-md overflow-x-auto w-full">
              <code>{JSON.stringify(responseData(), null, 2)}</code>
            </pre>
          </div>
        </Show>
      </div>
    </div>
  );
}
