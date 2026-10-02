import { Header } from "@/components/Header/Header";
import { HomeHero } from "@/components/HomeHero/HomeHero";

// The home page will act as the marketing site
function Home() {
  return <main className="flex flex-1 flex-col gap-6 p-8 md:px-20">
    <Header />
    <HomeHero />
  </main>;
}

export default Home;
