import Developer from "~/components/homePage/Developer";
import Features from "~/components/homePage/Features";
import Footer from "~/components/homePage/Footer";
import Hero from "~/components/homePage/Hero";
import Navbar from "~/components/homePage/Navbar";
import Technologies from "~/components/homePage/Technologies";

export default function Home() {
  return (
    <div class="min-h-screen bg-gradient-to-br from-gray-900 via-blue-900 to-blue-950 text-white flex flex-col">
      <Navbar />
      <Hero />
      <Features />
      <Technologies />
      <Developer />
      <Footer />
    </div>
  );
}
