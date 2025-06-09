import { Link, Meta, Title } from "@solidjs/meta";
import Developer from "~/components/homePage/Developer";
import Features from "~/components/homePage/Features";
import Footer from "~/components/homePage/Footer";
import Hero from "~/components/homePage/Hero";
import Navbar from "~/components/homePage/Navbar";
import Technologies from "~/components/homePage/Technologies";

export default function Home() {
  return (
    <div class="min-h-screen bg-gradient-to-br from-gray-900 via-blue-900 to-blue-950 text-white flex flex-col">
      <Title>Blue Fox - Modern Open-Source Real-time Chat Application</Title>
      <Meta
        name="description"
        content="Blue Fox is a blazing fast, open-source real-time chat application built with SolidJS and Tailwind CSS. Connect instantly and join a vibrant community."
      />
      <Meta
        name="keywords"
        content="Blue Fox, chat app, real-time chat, open source, SolidJS, Tailwind CSS, community, communication, web app, real-time messaging platform, SolidJS web application example"
      />
      <Link rel="canonical" href="https://blue-fox.vercel.app/" />
      <Navbar />
      <Hero />
      <Features />
      <Technologies />
      <Developer />
      <Footer />
    </div>
  );
}
