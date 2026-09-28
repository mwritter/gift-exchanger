
import { Gift } from "lucide-react";
import { Button } from "../ui/button";
import Link from "next/link";

export function Header() {
    return <div className="flex items-center w-full">
        <div className="flex-1 flex gap-2 items-center">
            <Gift />
            <p className="text-2xl font-bold">GiftExchanger</p>
        </div>
        <Button nativeButton={false} render={<Link href="/login" />}>Log in</Button>
    </div>
}