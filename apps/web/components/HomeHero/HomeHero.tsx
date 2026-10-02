"use client"

import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";

export function HomeHero() {
  const { push } = useRouter()

  return (
    <div className="grid md:grid-cols-2 gap-4 flex-1 p-5">
      <div className="flex flex-col gap-4 mx-auto max-w-125 justify-center text-center">
        <h1 className="font-bold md:text-5xl text-2xl">
          Make gift giving more meaningful.
        </h1>
        <p className="text-muted-foreground text-sm">Create a gift exchange, invite your friends, share wishlists, and let GiftExchanger handle the rest - including the secret assignments.</p>
        <div className="flex justify-center gap-2 w-full">
          <Button size="lg" onClick={() => push("/login")}>Create your first exchange</Button>
          <Button size="lg" className="text-primary" variant='outline'>Learn more</Button>
        </div>
      </div>
      <div></div>
    </div>
  );
}
