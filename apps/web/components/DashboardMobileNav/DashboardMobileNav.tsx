import { DashboardMobileNavLink } from "./DashboardMobileNavLink";
import navLinks from "../consts/navLinks";

export function DashboardMobileNav() {
    return <div className="md:hidden fixed bottom-0 flex justify-between w-full text-sm border-t-2 bg-background">
        {navLinks.map(({ href, Icon, text }) => <DashboardMobileNavLink key={href} href={href}>
            <Icon size={20} />
            {text}
        </DashboardMobileNavLink>)}
    </div>
}