"use client"

import { useState } from "react";
import { ExchangeInviteInput } from "./ExchangeInviteInput";
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { Trash } from "lucide-react";


export function ExchangeInviteForm() {
    const [emails, setEmails] = useState<string[]>([])

    const onEmailInputSubmit = (email: string) => {
        setEmails(cur => {
            const found = cur.find(e => e === email)
            if (!found) return [...cur, email]

            return cur
        })
    }

    const onEmailInputRemove = (email: string) => {
        setEmails(cur => cur.filter(e => e !== email))
    }


    return <div>
        <ul className="flex flex-col gap-2">
            <li>
                <ExchangeInviteInput onSubmit={onEmailInputSubmit} />
            </li>
            {emails.map((e, idx) => <li key={idx} className="flex gap-2 justify-center items-center">
                <Input value={e} disabled />
                <Button type="button" className="" variant='destructive' onClick={() => onEmailInputRemove(e)}>
                    <Trash />
                </Button>
            </li>)}
        </ul>
    </div>
}