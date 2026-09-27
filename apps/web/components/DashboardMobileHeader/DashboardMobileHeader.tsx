import { Gift } from "lucide-react";
import { LogoutButton } from "../LogoutButton/LogoutButton";

export function DashboardMobileHeader() {
    return <div className="md:hidden flex justify-between px-5 py-2 border-b-2">
        <div className="flex gap-2 items-center">
            <Gift size={20} className="text-primary" />
            <p className="font-bold text-md">GiftExchanger</p>
        </div>
        <LogoutButton />
    </div>
}