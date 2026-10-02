import { DashboardPageContent } from "@/components/DashboardPageLayout/DashboardPageContent";
import { DashboardPageHeader } from "@/components/DashboardPageLayout/DashboardPageHeader";

export default function WishlistsPage() {
    return (
        <>
            <DashboardPageHeader title="Wishlist" description="A wishlist of items you'd like to get! Add items to exchanges and your wish might get picked!" />
            <DashboardPageContent>
                Wishlist Page
            </DashboardPageContent>
        </>
    )
}