"use client"

import { ExchangeInviteInput } from "./ExchangeInviteInput";
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { Trash } from "lucide-react";

type Props = {
    emails: string[]
    onChange: (emails: string[]) => void
}

export function ExchangeInviteForm({ emails, onChange }: Props) {
    const onEmailInputSubmit = (email: string) => {
        if (emails.includes(email)) return
        onChange([...emails, email])
    }

    const onEmailInputRemove = (email: string) => {
        onChange(emails.filter(e => e !== email))
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