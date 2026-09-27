"use client"

import { DashboardPageContent } from "@/components/DashboardPageLayout/DashboardPageContent";
import { DashboardPageHeader } from "@/components/DashboardPageLayout/DashboardPageHeader";
import { DashboardPageLayout } from "@/components/DashboardPageLayout/DashboardPageLayout";
import { Button } from "@/components/ui/button";
import { Plus, X } from "lucide-react";
import { useState } from "react";

// List exchanges
// Create and update exchanges

export default function ExchangesPage() {
    const [showCreateNewForm, setShowCreateNewForm] = useState(false)

    return <DashboardPageLayout>
        {/* <div className="flex items-center justify-between md:max-w-125 gap-5">
            <DashboardPageHeader title="Exchanges" description="This is a list of exchanges your current a part of!" />
            {!showCreateNewForm ? <Button onClick={() => setShowCreateNewForm(true)} size='sm'>
                <Plus />
                Create
            </Button> : <Button variant={'destructive'} onClick={() => setShowCreateNewForm(false)}><X /></Button>}
        </div> */}
        <DashboardPageContent>
            <EmptyExchangeList />
        </DashboardPageContent>
    </DashboardPageLayout>
}

function EmptyExchangeList() {
    return <div className="w-full text-center justify-center">
        <p className="text-sm">No exchanges to list yet</p>
    </div>
}