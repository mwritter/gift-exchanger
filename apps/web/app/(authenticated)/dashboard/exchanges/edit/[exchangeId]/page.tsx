"use client"

import { ExchangeForm } from "@/components/ExchangeForm/ExchangeForm"
import { use } from "react"

export default function EditExchange({
    params,
}: {
    params: Promise<{ exchangeId: string }>
}) {
    const { exchangeId } = use(params)

    // Fetch the exchange to fill in the form

    const handleSubmit = () => {
        console.log('submitting edit ' + exchangeId)
    }

    return <div>
        {/* Pass exchange values here */}
        <ExchangeForm title="Edit" onSubmit={handleSubmit} />
    </div>
}