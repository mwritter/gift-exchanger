import { Gift, Heart, Home, LucideIcon, Settings } from "lucide-react";

export type NavLink = {
  href: string;
  Icon: LucideIcon;
  text: string;
};

export default [
  {
    href: "/dashboard",
    Icon: Home,
    text: "Home",
  },
  {
    href: "/dashboard/exchanges",
    Icon: Gift,
    text: "Exchanges",
  },
  {
    href: "/dashboard/wishlist",
    Icon: Heart,
    text: "Wishlist",
  },
  {
    href: "/dashboard/settings",
    Icon: Settings,
    text: "Settings",
  },
] as const satisfies NavLink[];
