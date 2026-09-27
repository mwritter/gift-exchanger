"use client"

import { ExchangeForm } from "@/components/ExchangeForm/ExchangeForm";

export default function CreateExchangePage() {

    const handleSubmit = () => {
        console.log("submitting create")
    }

    return <div>
        <ExchangeForm title="Create" onSubmit={handleSubmit} />
    </div>
}