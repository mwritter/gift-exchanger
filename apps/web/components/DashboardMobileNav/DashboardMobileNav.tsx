import { Gift, Heart, Home, Settings } from "lucide-react";
import { DashboardMobileNavLink } from "./DashboardMobileNavLink";

// TODO: Clean this up

export function DashboardMobileNav() {
    return <div className="md:hidden fixed bottom-0 flex justify-between w-full text-sm border-t-2 bg-background">
        <DashboardMobileNavLink href="/dashboard">
            <Home size={20} />
            Home
        </DashboardMobileNavLink>
        <DashboardMobileNavLink href="/dashboard/exchanges">
            <Gift size={20} />
            Exchanges
        </DashboardMobileNavLink>
        <DashboardMobileNavLink href="/dashboard/wishlist">
            <Heart size={20} />
            Wishlist
        </DashboardMobileNavLink>
        <DashboardMobileNavLink href="/dashboard/settings">
            <Settings size={20} />
            Settings
        </DashboardMobileNavLink>
    </div>
}