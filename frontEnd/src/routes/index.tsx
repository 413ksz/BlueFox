import Developer from "~/components/Developer";
import Features from "~/components/Features";
import Footer from "~/components/Footer";
import Hero from "~/components/Hero";
import Navbar from "~/components/Navbar";
import Technologies from "~/components/Technologies";

export default function Home() {
  return (
    <div class="min-h-screen bg-gradient-to-br from-gray-900 via-blue-900 to-blue-950 text-white flex flex-col">
      <Navbar/>
      <Hero/>
      <Features/>
      <Technologies/>
      <Developer/>
      <Footer/>
    </div>
  );
}
